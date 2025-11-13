package repository

import (
	"context"
	"fmt"
	"log/slog"
	"pubsub-ckg-tb/internal/config"
	"pubsub-ckg-tb/internal/db/connection"
	"pubsub-ckg-tb/internal/db/dbtypes"
	"pubsub-ckg-tb/internal/db/utils"
	"pubsub-ckg-tb/internal/models"
	"slices"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CKGTB interface {
	GetPendingTbSkrining(start string, end string, limit int64) ([]models.SkriningCKGResult, error)
	GetOnePendingTbSkrining(table string, docBytes []byte) (*models.SkriningCKGResult, error)
	UpdateTbPatientStatus(input []models.StatusPasien) ([]models.StatusPasienResult, error)
}

type CKGTBRepository struct {
	Configurations *config.Configurations
	Connnection    connection.DatabaseConnection
	Context        context.Context
	// MasterWilayahTable any
	// MasterFaskesTable  any
	// TransactionTable   any
	// TbPatientTable     any

	useCache     bool
	cacheWilayah map[string]models.MasterWilayah
	cacheFaskes  map[string]models.MasterFaskes
}

func NewCKGTBRepository(ctx context.Context, config *config.Configurations, conn connection.DatabaseConnection) *CKGTBRepository {

	// _, masterWilayahTable := utils.GetCollection(ctx, conn, config.CKG.TableMasterWilayah, 0)
	// _, masterFaskesTable := utils.GetCollection(ctx, conn, config.CKG.TableMasterFaskes, 0)
	// _, tbPatientTable := utils.GetCollection(ctx, conn, config.CKG.TableStatus, 0)
	// _, transactionTable := utils.GetCollection(ctx, conn, config.CKG.TableSkrining, 0)

	return &CKGTBRepository{
		Configurations: config,
		Connnection:    conn,
		Context:        ctx,
		useCache:       config.CKG.UseCache,
		// MasterWilayahTable: masterWilayahTable,
		// MasterFaskesTable:  masterFaskesTable,
		// TbPatientTable:     tbPatientTable,
		// TransactionTable:   transactionTable,
	}
}

func (r *CKGTBRepository) GetPendingTbSkrining(start string, end string, limit int64) ([]models.SkriningCKGResult, error) {
	// Get Skrining
	filter := dbtypes.M{
		"updated_at": dbtypes.M{
			"$gte": start,
			"$lte": end,
		},
	}
	ret, err := r.Connnection.Find(r.Context, r.Configurations.CKG.TableSkrining, nil, filter, nil, limit, 0)
	if err != nil {
		slog.Debug("GetPendingTbSkrining:", "error", err)
		return nil, err
	}
	result := []models.SkriningCKGResult{}
	for _, entry := range ret.([]dbtypes.M) {
		raw := models.SkriningCKGRaw{}
		raw.FromMap(entry)
		res := raw.ToSkriningCKGResult()
		r._HitungHasilSkrining(raw, &res)
		r._MappingMasterData(r.Context, r.Context, raw, &res)
		result = append(result, res)
	}

	return result, nil
}

func (r *CKGTBRepository) GetOnePendingTbSkrining(id string, docBytes []byte) (*models.SkriningCKGResult, error) {
	// Convert to SkriningCKGRaw
	var raw models.SkriningCKGRaw

	if err := bson.Unmarshal(docBytes, &raw); err != nil {
		return nil, fmt.Errorf("gagal unmarshal document: %v", err)
	}

	// Create new SkriningCKGResult from raw data
	res := raw.ToSkriningCKGResult()
	/*res := &models.SkriningCKGResult{
		PasienCKGID:                raw.PasienCKGID,
		PasienNIK:                  raw.PasienNIK,
		PasienNama:                 raw.PasienNama,
		PasienJenisKelamin:         raw.PasienJenisKelamin,
		PasienTglLahir:             raw.PasienTglLahir,
		PasienUsia:                 raw.PasienUsia,
		PasienPekerjaan:            raw.PasienPekerjaan,
		PasienAlamat:               raw.PasienAlamat,
		PasienNoHandphone:          raw.PasienNoHandphone,
		TglPemeriksaan:             raw.TglPemeriksaan,
		BeratBadan:                 raw.BeratBadan,
		TinggiBadan:                raw.TinggiBadan,
		StatusImt:                  raw.StatusImt,
		HasilGds:                   raw.HasilGds,
		HasilGdp:                   raw.HasilGdp,
		HasilGdpp:                  raw.HasilGdpp,
		KekuranganGizi:             raw.KekuranganGizi,
		Merokok:                    raw.Merokok,
		PerokokPasif:               raw.PerokokPasif,
		LansiaDiatas65:             raw.LansiaDiatas65,
		IbuHamil:                   raw.IbuHamil,
		RiwayatDm:                  raw.HasilPemeriksaanDm,
		RiwayatHt:                  raw.HasilPemeriksaanHt,
		InfeksiHivAids:             raw.InfeksiHivAids,
		GejalaBatuk:                raw.GejalaDanTandaBatuk,
		GejalaBbTurun:              raw.GejalaDanTandaBbTurun,
		GejalaDemamHilangTimbul:    raw.GejalaDanTandaDemamHilangTimbul,
		GejalaLesuMalaise:          raw.GejalaDanTandaLesuMalaise,
		GejalaBerkeringatMalam:     raw.GejalaDanTandaBerkeringatMalam,
		GejalaPembesaranKelenjarGB: raw.GejalaDanTandaPembesaranKelenjarGB,
		KontakPasienTbc:            raw.KontakPasienTbc,
		HasilPemeriksaanTbBta:      raw.HasilPemeriksaanTbBta,
		HasilPemeriksaanTbTcm:      raw.HasilPemeriksaanTbTcm,
		HasilPemeriksaanPoct:       raw.HasilPemeriksaanPoct,
		HasilPemeriksaanRadiologi:  raw.HasilPemeriksaanRadiologi,
	}*/

	r._HitungHasilSkrining(raw, &res)
	r._MappingMasterData(r.Context, r.Context, raw, &res)

	return &res, nil
}

func (r *CKGTBRepository) UpdateTbPatientStatus(input []models.StatusPasien) ([]models.StatusPasienResult, error) {
	results := make([]models.StatusPasienResult, 0, len(input))
	collectionName := r.Configurations.CKG.TableStatus

	for i, item := range input {
		res := models.StatusPasienResult{
			PasienCkgID: item.PasienCkgID,
			TerdugaID:   item.TerdugaID,
			PasienTbID:  item.PasienTbID,
			PasienNIK:   item.PasienNIK,
			IsError:     false,
			Respons:     "",
		}

		// Validas data input
		err := r._ValidateSkriningData(item, i)
		if err != nil {
			res.IsError = true
			res.Respons = err.Error()
			results = append(results, res)
			continue
		}

		// Simpan atau update database.
		resExist, err := utils.FindPasienTb(r.Context, r.Connnection, collectionName, item)
		if resExist != nil { // sudah ada status
			if utils.IsNotEmptyString(resExist.PasienCkgID) {
				res.PasienCkgID = resExist.PasienCkgID
			}
			// Ditemukan tapi id CKG belum di set (kemungkinan SITB mendahului lapor)
			if !utils.IsNotEmptyString(resExist.PasienCkgID) && utils.IsNotEmptyString(item.PasienNIK) {
				var transaction models.SkriningCKGRaw
				// filterTx := bson.D{
				// 	{Key: "nik", Value: item.PasienNIK},
				// }
				filterTx := dbtypes.M{
					"nik": item.PasienNIK,
				}

				// errTx := r.TransactionTable.FindOne(r.Context, filterTx).Decode(&transaction)
				errTx := r.Connnection.FindOne(r.Context, &transaction, collectionName, nil, filterTx, nil)
				if errTx == nil { // Transaksi layanan CKG ditemukan
					item.PasienCkgID = &transaction.PasienCKGID
					res.PasienCkgID = &transaction.PasienCKGID
				}
			}

			// jaga-jaga kalau pasienTbID tidak dikirim oleh sitb di pengiriman berikutnya
			if utils.IsNotEmptyString(resExist.PasienTbID) && item.PasienTbID == nil {
				item.PasienTbID = resExist.PasienTbID
				res.PasienTbID = resExist.PasienTbID
			}

			if utils.IsNotEmptyString(item.PasienTbID) {
				if utils.IsNotEmptyString(item.StatusDiagnosis) {
					if !utils.IsNotEmptyString(item.DiagnosisLabHasilTCM) {
						item.DiagnosisLabHasilTCM = nil
					}
					if !utils.IsNotEmptyString(item.DiagnosisLabHasilBTA) {
						item.DiagnosisLabHasilBTA = nil
					}
				} else {
					item.StatusDiagnosis = nil
					item.DiagnosisLabHasilTCM = nil
					item.DiagnosisLabHasilBTA = nil
					item.TanggalMulaiPengobatan = nil
					item.TanggalSelesaiPengobatan = nil
					item.HasilAkhir = nil
				}
			} else {
				item.StatusDiagnosis = nil
				item.DiagnosisLabHasilTCM = nil
				item.DiagnosisLabHasilBTA = nil
				item.TanggalMulaiPengobatan = nil
				item.TanggalSelesaiPengobatan = nil
				item.HasilAkhir = nil
			}

			msg, err1 := utils.UpdatePasienTb(r.Context, r.Connnection, collectionName, item)
			res.Respons = msg
			if err1 != nil {
				res.IsError = true
			}
		} else if err == mongo.ErrNoDocuments { // status baru
			// Coba cari di transaksi
			if utils.IsNotEmptyString(item.PasienNIK) {
				var transaction models.SkriningCKGRaw
				// filterTx := bson.D{
				// 	{Key: "nik", Value: item.PasienNIK},
				// }
				filterTx := dbtypes.M{
					"nik": item.PasienNIK,
				}

				// errTx := r.TransactionTable.FindOne(r.Context, filterTx).Decode(&transaction)
				errTx := r.Connnection.FindOne(r.Context, &transaction, collectionName, nil, filterTx, nil)
				if errTx == nil { // Transaksi layanan CKG ditemukan
					item.PasienCkgID = &transaction.PasienCKGID
					res.PasienCkgID = &transaction.PasienCKGID
				} else { // SITB duluan dilaporkan oleh CKG
					item.PasienCkgID = nil
					res.PasienCkgID = nil
				}
			}

			if utils.IsNotEmptyString(item.PasienTbID) {
				if utils.IsNotEmptyString(item.StatusDiagnosis) {
					if !utils.IsNotEmptyString(item.DiagnosisLabHasilTCM) {
						item.DiagnosisLabHasilTCM = nil
					}
					if !utils.IsNotEmptyString(item.DiagnosisLabHasilBTA) {
						item.DiagnosisLabHasilBTA = nil
					}
				} else {
					item.StatusDiagnosis = nil
					item.DiagnosisLabHasilTCM = nil
					item.DiagnosisLabHasilBTA = nil
					item.TanggalMulaiPengobatan = nil
					item.TanggalSelesaiPengobatan = nil
					item.HasilAkhir = nil
				}
			} else {
				item.StatusDiagnosis = nil
				item.DiagnosisLabHasilTCM = nil
				item.DiagnosisLabHasilBTA = nil
				item.TanggalMulaiPengobatan = nil
				item.TanggalSelesaiPengobatan = nil
				item.HasilAkhir = nil
			}

			msg, err1 := utils.InsertPasienTb(r.Context, r.Connnection, collectionName, item)
			res.Respons = msg
			if err1 != nil {
				res.IsError = true
			}
		} else {
			continue
		}

		results = append(results, res)
	}

	return results, nil
}

func (r *CKGTBRepository) _MappingMasterData(ctxMasterWilayah context.Context, ctxMasterFaskes context.Context, raw models.SkriningCKGRaw, res *models.SkriningCKGResult) {
	collectionNameMasterWilayah := r.Configurations.CKG.TableMasterWilayah
	if utils.IsNotEmptyString(raw.PasienKelurahan) {
		kelurahan, _ := utils.FindMasterWilayah(*raw.PasienKelurahan, ctxMasterWilayah, r.Connnection, collectionNameMasterWilayah, r.useCache, &r.cacheWilayah)
		if kelurahan != nil {
			res.PasienKelurahanSatusehat = raw.PasienKelurahan
			res.PasienKelurahanSitb = kelurahan.KelurahanID
		}
	}

	if utils.IsNotEmptyString(raw.PasienKecamatan) {
		kecamatan, _ := utils.FindMasterWilayah(*raw.PasienKecamatan, ctxMasterWilayah, r.Connnection, collectionNameMasterWilayah, r.useCache, &r.cacheWilayah)
		if kecamatan != nil {
			res.PasienKecamatanSatusehat = raw.PasienKecamatan
			res.PasienKecamatanSitb = kecamatan.KecamatanID
		}
	}

	if utils.IsNotEmptyString(raw.PasienKabkota) {
		kabupaten, _ := utils.FindMasterWilayah(*raw.PasienKabkota, ctxMasterWilayah, r.Connnection, collectionNameMasterWilayah, r.useCache, &r.cacheWilayah)
		if kabupaten != nil {
			res.PasienKabkotaSatusehat = raw.PasienKabkota
			res.PasienKabkotaSitb = kabupaten.KabupatenID
		}
	}

	if utils.IsNotEmptyString(raw.PasienProvinsi) {
		provinsi, _ := utils.FindMasterWilayah(*raw.PasienProvinsi, ctxMasterWilayah, r.Connnection, collectionNameMasterWilayah, r.useCache, &r.cacheWilayah)
		if provinsi != nil {
			res.PasienProvinsiSatusehat = raw.PasienProvinsi
			res.PasienProvinsiSitb = provinsi.ProvinsiID
		}
	}

	if utils.IsNotEmptyString(raw.KodeFaskes) {
		collectionNameMasterFaskes := r.Configurations.CKG.TableMasterFaskes
		faskes, _ := utils.FindMasterFaskes(*raw.KodeFaskes, ctxMasterFaskes, r.Connnection, collectionNameMasterFaskes, r.useCache, &r.cacheFaskes)
		if faskes != nil {
			res.KodeFaskesSatusehat = raw.KodeFaskes
			res.KodeFaskesSITB = faskes.ID
		}
	}
}

func (r *CKGTBRepository) _ValidateSkriningData(item models.StatusPasien, i int) error {
	// TerdugaID dan PasienNIK tidak boleh kosong
	if item.TerdugaID == nil || *item.TerdugaID == "" {
		return fmt.Errorf("validation error at index %d: terduga_id is required", i)
	}
	if item.PasienNIK == nil || *item.PasienNIK == "" {
		return fmt.Errorf("validation error at index %d: pasien_nik is required", i)
	}

	// 1=TBC SO,
	// 2=TBC RO,
	// 3= Bukan TBC
	statusDiagnosis := []string{"TBC SO", "TBC RO", "Bukan TBC"}

	// Paling tidak StatusTerduga, atau DiagnosisLabHasil harus ada
	if item.PasienTbID != nil && (item.StatusDiagnosis == nil || !slices.Contains(statusDiagnosis, *item.StatusDiagnosis)) {
		return fmt.Errorf("validation error at index %d: at least one of status_terduga, or status_diagnosa must be provided", i)
	} else {
		item.StatusDiagnosis = nil // abaikan
	}

	if utils.IsNotEmptyString(item.StatusDiagnosis) {
		if !utils.IsNotEmptyString(item.DiagnosisLabHasilTCM) {
			return fmt.Errorf("validation error at index %d: diagnosis_lab_hasil_tcm is required when status_diagnosa is provided", i)
		}
		if !utils.IsNotEmptyString(item.DiagnosisLabHasilBTA) {
			return fmt.Errorf("validation error at index %d: diagnosis_lab_hasil_bta is required when status_diagnosa is provided", i)
		}
	}

	// Hasil Akhir
	// 1= Sembuh,
	// 2= Pengobatan Lengkap,
	// 3= Pengobatan Gagal ,
	// 4= Meninggal,
	// 5= Putus berobat (lost to follow up),
	// 6= Tidak dievaluasi/pindah,
	// 7= Gagal karena Perubahan Diagnosis, "
	statusAkhir := []string{"Sembuh", "Pengobatan Lengkap", "Pengobatan Gagal", "Meninggal", "Putus berobat (lost to follow up)", "Tidak dievaluasi/pindah", "Gagal karena Perubahan Diagnosis"}
	if item.HasilAkhir != nil && !slices.Contains(statusAkhir, *item.HasilAkhir) {
		return fmt.Errorf("validation error at index %d: hasil_akhir must be one of %v", i, statusAkhir)
	}

	return nil
}

func (r *CKGTBRepository) _HitungHasilSkrining(raw models.SkriningCKGRaw, res *models.SkriningCKGResult) {
	hasilSkrining := "Tidak"

	if raw.PasienUsia < 15 {
		// Gejala batuk dan sudah lebih dari 14 hari
		if raw.GejalaDanTandaBatuk != nil && *raw.GejalaDanTandaBatuk == "Ya" {
			hasilSkrining = "Ya"
		}
		if raw.GejalaDanTandaBbTurun != nil && *raw.GejalaDanTandaBbTurun == "Ya" {
			hasilSkrining = "Ya"
		}
		if raw.GejalaDanTandaDemamHilangTimbul != nil && *raw.GejalaDanTandaDemamHilangTimbul == "Ya" {
			hasilSkrining = "Ya"
		}
		if raw.GejalaDanTandaLesuMalaise != nil && *raw.GejalaDanTandaLesuMalaise == "Ya" {
			hasilSkrining = "Ya"
		}

		// bersihkan gejala untuk dewasa
		res.GejalaBerkeringatMalam = nil
		res.GejalaPembesaranKelenjarGB = nil
	} else { // 15 tahun ke atas
		if raw.InfeksiHivAids != nil && *raw.InfeksiHivAids == "Ya" {
			// Cukup gejala batuk tanpa harus melihat sudah 14 hari atau tidak
			if raw.GejalaDanTandaBatuk != nil && *raw.GejalaDanTandaBatuk == "Ya" {
				hasilSkrining = "Ya"
			}
			if raw.GejalaDanTandaBbTurun != nil && *raw.GejalaDanTandaBbTurun == "Ya" {
				hasilSkrining = "Ya"
			}
			if raw.GejalaDanTandaDemamHilangTimbul != nil && *raw.GejalaDanTandaDemamHilangTimbul == "Ya" {
				hasilSkrining = "Ya"
			}
			if raw.GejalaDanTandaBerkeringatMalam != nil && *raw.GejalaDanTandaBerkeringatMalam == "Ya" {
				hasilSkrining = "Ya"
			}
			if raw.GejalaDanTandaPembesaranKelenjarGB != nil && *raw.GejalaDanTandaPembesaranKelenjarGB == "Ya" {
				hasilSkrining = "Ya"
			}
		} else {
			// Gejala batuk dan sudah lebih dari 14 hari
			if raw.GejalaDanTandaBatuk != nil && *raw.GejalaDanTandaBatuk == "Ya" {
				hasilSkrining = "Ya"
			}
			if raw.GejalaDanTandaBbTurun != nil && *raw.GejalaDanTandaBbTurun == "Ya" {
				hasilSkrining = "Ya"
			}
			if raw.GejalaDanTandaDemamHilangTimbul != nil && *raw.GejalaDanTandaDemamHilangTimbul == "Ya" {
				hasilSkrining = "Ya"
			}
			if raw.GejalaDanTandaBerkeringatMalam != nil && *raw.GejalaDanTandaBerkeringatMalam == "Ya" {
				hasilSkrining = "Ya"
			}
			if raw.GejalaDanTandaPembesaranKelenjarGB != nil && *raw.GejalaDanTandaPembesaranKelenjarGB == "Ya" {
				hasilSkrining = "Ya"
			}
		}

		// bersihkan gejala untuk anak
		res.GejalaLesuMalaise = nil
	}

	res.HasilSkriningTbc = &hasilSkrining
	if hasilSkrining == "Ya" {
		if raw.HasilPemeriksaanTbTcm != nil {
			//TODO: koordinasikan mapping nilai TCM dengan DE
			// convert hasil TCM ke ["not_detected", "rif_sen", "rif_res", "rif_indet", "invalid", "error", "no_result", "tdl"]
			mapTcm := map[string]string{
				"neg":       "not_detected",
				"rif sen":   "rif_sen",
				"rif res":   "rif_res",
				"rif indet": "rif_indet",
				"invalid":   "invalid",
				"error":     "error",
				"no result": "no_result",
			}
			if utils.IsNotEmptyString(raw.HasilPemeriksaanTbTcm) {
				tcm := strings.ToLower(*raw.HasilPemeriksaanTbTcm)
				if val, ok := mapTcm[tcm]; ok {
					res.HasilPemeriksaanTbTcm = &val
				}
			}
		}

		if raw.HasilPemeriksaanTbBta != nil {
			//TODO: koordinasikan mapping nilai BTA dengan DE
			// convert hasil BTA ke ["negatif", "positif"]
			var hasilTbBta *string
			if utils.IsNotEmptyString(raw.HasilPemeriksaanTbBta) {
				bta := strings.ToLower(*raw.HasilPemeriksaanTbBta)
				hasilTbBta = &bta
			}
			res.HasilPemeriksaanTbBta = hasilTbBta
		}

		if raw.HasilPemeriksaanPoct != nil {
			// convert hasil POCT ke ["negatif", "positif"]
			var hasilTbPoct *string
			if utils.IsNotEmptyString(raw.HasilPemeriksaanPoct) {
				poct := strings.ToLower(*raw.HasilPemeriksaanPoct)
				hasilTbPoct = &poct
			}
			res.HasilPemeriksaanPoct = hasilTbPoct
		}

		if raw.HasilPemeriksaanRadiologi != nil {
			// convert hasil Radiologi ke ["normal", "abnormalitas-tbc", "abnormalitas-bukan-tbc"]
			var hasilHasilPemeriksaanRadiologi *string
			if utils.IsNotEmptyString(raw.HasilPemeriksaanRadiologi) {
				radiologi := strings.ReplaceAll(strings.ToLower(*raw.HasilPemeriksaanRadiologi), " ", "-")
				hasilHasilPemeriksaanRadiologi = &radiologi
			}
			res.HasilPemeriksaanRadiologi = hasilHasilPemeriksaanRadiologi
		}

		res.TerdugaTb = &hasilSkrining
	}
}
