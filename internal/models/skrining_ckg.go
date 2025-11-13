package models

// Hasil Pemrosesan dari database
type SkriningCKGRaw struct {
	// Identitas Pasien
	PasienCKGID        string  `bson:"pasien_id"`
	PasienNIK          string  `bson:"nik"`
	PasienNama         string  `bson:"pasien_name"`
	PasienJenisKelamin string  `bson:"jenis_kelamin"`
	PasienTglLahir     string  `bson:"tgl_lahir"`
	PasienUsia         int     `bson:"usia"`
	PasienPekerjaan    *string `bson:"pekerjaan"` //TODO: saat ini data belum tersedia di dwh mongodb
	PasienProvinsi     *string `bson:"provinsi_pasien"`
	PasienKabkota      *string `bson:"kabkota_pasien"`
	PasienKecamatan    *string `bson:"kecamatan_pasien"`
	PasienKelurahan    *string `bson:"kelurahan_pasien"`
	PasienAlamat       *string `bson:"alamat"`
	PasienNoHandphone  string  `bson:"no_handphone"`

	// Data Kunjungan
	KodeFaskes     *string `bson:"kode_faskes"`
	NamaFaskes     *string `bson:"nama_faskes"`
	ProvinsiFaskes *string `bson:"provinsi_faskes"`
	KabkotaFaskes  *string `bson:"kabkota_faskes"`
	TglPemeriksaan string  `bson:"tgl_pemeriksaan"`

	// Data Hasil Pemeriksaan
	BeratBadan                *float64 `bson:"berat_badan"`
	TinggiBadan               *float64 `bson:"tinggi_badan"`
	StatusImt                 *string  `bson:"imt"`
	KekuranganGizi            *string  `bson:"kekurangan_gizi"`
	Merokok                   *string  `bson:"merokok"`
	PerokokPasif              *string  `bson:"perokok_pasif"`
	LansiaDiatas65            *string  `bson:"lansia_lebih_dari_65"`
	IbuHamil                  *string  `bson:"ibu_hamil"`
	HasilGds                  *float64 `bson:"hasil_gds"`
	HasilGdp                  *float64 `bson:"hasil_gdp"`
	HasilGdpp                 *float64 `bson:"hasil_gdpp"`
	PemeriksaanChestXray      *string  `bson:"pemeriksaan_chest_xray"`
	HasilPemeriksaanTbBta     *string  `bson:"hasil_pemeriksaan_tb_bta"`
	HasilPemeriksaanTbTcm     *string  `bson:"hasil_pemeriksaan_tb_tcm"`
	HasilPemeriksaanDm        *string  `bson:"hasil_pemeriksaan_dm"`
	HasilPemeriksaanHt        *string  `bson:"hasil_pemeriksaan_ht"`
	HasilPemeriksaanPoct      *string  `bson:"hasil_pemeriksaan_tb_poct"`
	HasilPemeriksaanRadiologi *string  `bson:"hasil_pemeriksaan_tb_radiologi"`
	InfeksiHivAids            *string  `bson:"inveksi_hiv_aids"`

	// Data Skrining TB
	GejalaDanTandaBatuk                *string `bson:"gejala_dan_tanda_batuk"`
	GejalaDanTandaBbTurun              *string `bson:"gejala_dan_tanda_bb_turun"`                // dipersiapkan untuk variabel baru/tambahan Andak & Dewasa
	GejalaDanTandaDemamHilangTimbul    *string `bson:"gejala_dan_tanda_demam_hilang_timbul"`     // dipersiapkan untuk variabel baru/tambahan Andak & Dewasa
	GejalaDanTandaLesuMalaise          *string `bson:"gejala_dan_tanda_lesu_malaise"`            // dipersiapkan untuk variabel baru/tambahan Anak
	GejalaDanTandaBerkeringatMalam     *string `bson:"gejala_dan_tanda_berkeringat_malam"`       // dipersiapkan untuk variabel baru/tambahan Dewasa
	GejalaDanTandaPembesaranKelenjarGB *string `bson:"gejala_dan_tanda_pembesaran_getah_bening"` // dipersiapkan untuk variabel baru/tambahan Dewasa
	KontakPasienTbc                    *string `bson:"kontak_pasien_tbc"`
	GejalaDanTandaTbc                  *string `bson:"gejala_dan_tanda_tbc"`

	TindakLanjutPenegakanDiagnosa *string `bson:"tindak_lanjut_penegakan_diagnosa"`
	UpdatedAt                     string  `bson:"updated_at"`
}

// FromMap creates a SkriningCKGRaw from a map
func (s *SkriningCKGRaw) FromMap(data map[string]any) {
	// Map fields
	if val, ok := data["pasien_id"].(string); ok {
		s.PasienCKGID = val
	}
	if val, ok := data["nik"].(string); ok {
		s.PasienNIK = val
	}
	if val, ok := data["pasien_name"].(string); ok {
		s.PasienNama = val
	}
	if val, ok := data["jenis_kelamin"].(string); ok {
		s.PasienJenisKelamin = val
	}
	if val, ok := data["tgl_lahir"].(string); ok {
		s.PasienTglLahir = val
	}
	if val, ok := data["usia"].(float64); ok {
		s.PasienUsia = int(val)
	}
	if val, ok := data["pekerjaan"].(string); ok {
		s.PasienPekerjaan = &val
	}
	if val, ok := data["provinsi_pasien"].(string); ok {
		s.PasienProvinsi = &val
	}
	if val, ok := data["kabkota_pasien"].(string); ok {
		s.PasienKabkota = &val
	}
	if val, ok := data["kecamatan_pasien"].(string); ok {
		s.PasienKecamatan = &val
	}
	if val, ok := data["kelurahan_pasien"].(string); ok {
		s.PasienKelurahan = &val
	}
	if val, ok := data["alamat"].(string); ok {
		s.PasienAlamat = &val
	}
	if val, ok := data["no_handphone"].(string); ok {
		s.PasienNoHandphone = val
	}
	if val, ok := data["kode_faskes"].(string); ok {
		s.KodeFaskes = &val
	}
	if val, ok := data["nama_faskes"].(string); ok {
		s.NamaFaskes = &val
	}
	if val, ok := data["provinsi_faskes"].(string); ok {
		s.ProvinsiFaskes = &val
	}
	if val, ok := data["kabkota_faskes"].(string); ok {
		s.KabkotaFaskes = &val
	}
	if val, ok := data["tgl_pemeriksaan"].(string); ok {
		s.TglPemeriksaan = val
	}
	if val, ok := data["berat_badan"].(float64); ok {
		s.BeratBadan = &val
	}
	if val, ok := data["tinggi_badan"].(float64); ok {
		s.TinggiBadan = &val
	}
	if val, ok := data["imt"].(string); ok {
		s.StatusImt = &val
	}
	if val, ok := data["kekurangan_gizi"].(string); ok {
		s.KekuranganGizi = &val
	}
	if val, ok := data["merokok"].(string); ok {
		s.Merokok = &val
	}
	if val, ok := data["perokok_pasif"].(string); ok {
		s.PerokokPasif = &val
	}
	if val, ok := data["lansia_lebih_dari_65"].(string); ok {
		s.LansiaDiatas65 = &val
	}
	if val, ok := data["ibu_hamil"].(string); ok {
		s.IbuHamil = &val
	}
	if val, ok := data["hasil_gds"].(float64); ok {
		s.HasilGds = &val
	}
	if val, ok := data["hasil_gdp"].(float64); ok {
		s.HasilGdp = &val
	}
	if val, ok := data["hasil_gdpp"].(float64); ok {
		s.HasilGdpp = &val
	}
	if val, ok := data["pemeriksaan_chest_xray"].(string); ok {
		s.PemeriksaanChestXray = &val
	}
	if val, ok := data["hasil_pemeriksaan_tb_bta"].(string); ok {
		s.HasilPemeriksaanTbBta = &val
	}
	if val, ok := data["hasil_pemeriksaan_tb_tcm"].(string); ok {
		s.HasilPemeriksaanTbTcm = &val
	}
	if val, ok := data["hasil_pemeriksaan_dm"].(string); ok {
		s.HasilPemeriksaanDm = &val
	}
	if val, ok := data["hasil_pemeriksaan_ht"].(string); ok {
		s.HasilPemeriksaanHt = &val
	}
	if val, ok := data["hasil_pemeriksaan_tb_poct"].(string); ok {
		s.HasilPemeriksaanPoct = &val
	}
	if val, ok := data["hasil_pemeriksaan_tb_radiologi"].(string); ok {
		s.HasilPemeriksaanRadiologi = &val
	}
	if val, ok := data["inveksi_hiv_aids"].(string); ok {
		s.InfeksiHivAids = &val
	}
	if val, ok := data["gejala_dan_tanda_batuk"].(string); ok {
		s.GejalaDanTandaBatuk = &val
	}
	if val, ok := data["gejala_dan_tanda_bb_turun"].(string); ok {
		s.GejalaDanTandaBbTurun = &val
	}
	if val, ok := data["gejala_dan_tanda_demam_hilang_timbul"].(string); ok {
		s.GejalaDanTandaDemamHilangTimbul = &val
	}
	if val, ok := data["gejala_dan_tanda_lesu_malaise"].(string); ok {
		s.GejalaDanTandaLesuMalaise = &val
	}
	if val, ok := data["gejala_dan_tanda_berkeringat_malam"].(string); ok {
		s.GejalaDanTandaBerkeringatMalam = &val
	}
	if val, ok := data["gejala_dan_tanda_pembesaran_getah_bening"].(string); ok {
		s.GejalaDanTandaPembesaranKelenjarGB = &val
	}
	if val, ok := data["kontak_pasien_tbc"].(string); ok {
		s.KontakPasienTbc = &val
	}
	if val, ok := data["gejala_dan_tanda_tbc"].(string); ok {
		s.GejalaDanTandaTbc = &val
	}
	if val, ok := data["tindak_lanjut_penegakan_diagnosa"].(string); ok {
		s.TindakLanjutPenegakanDiagnosa = &val
	}
	if val, ok := data["updated_at"].(string); ok {
		s.UpdatedAt = val
	}
}

// ToMap converts SkriningCKGRaw to a map
func (s *SkriningCKGRaw) ToMap() map[string]any {
	return map[string]any{
		"pasien_id":                            s.PasienCKGID,
		"nik":                                  s.PasienNIK,
		"pasien_name":                          s.PasienNama,
		"jenis_kelamin":                        s.PasienJenisKelamin,
		"tgl_lahir":                            s.PasienTglLahir,
		"usia":                                 s.PasienUsia,
		"pekerjaan":                            s.PasienPekerjaan,
		"provinsi_pasien":                      s.PasienProvinsi,
		"kabkota_pasien":                       s.PasienKabkota,
		"kecamatan_pasien":                     s.PasienKecamatan,
		"kelurahan_pasien":                     s.PasienKelurahan,
		"alamat":                               s.PasienAlamat,
		"no_handphone":                         s.PasienNoHandphone,
		"kode_faskes":                          s.KodeFaskes,
		"nama_faskes":                          s.NamaFaskes,
		"provinsi_faskes":                      s.ProvinsiFaskes,
		"kabkota_faskes":                       s.KabkotaFaskes,
		"tgl_pemeriksaan":                      s.TglPemeriksaan,
		"berat_badan":                          s.BeratBadan,
		"tinggi_badan":                         s.TinggiBadan,
		"imt":                                  s.StatusImt,
		"kekurangan_gizi":                      s.KekuranganGizi,
		"merokok":                              s.Merokok,
		"perokok_pasif":                        s.PerokokPasif,
		"lansia_lebih_dari_65":                 s.LansiaDiatas65,
		"ibu_hamil":                            s.IbuHamil,
		"hasil_gds":                            s.HasilGds,
		"hasil_gdp":                            s.HasilGdp,
		"hasil_gdpp":                           s.HasilGdpp,
		"pemeriksaan_chest_xray":               s.PemeriksaanChestXray,
		"hasil_pemeriksaan_tb_bta":             s.HasilPemeriksaanTbBta,
		"hasil_pemeriksaan_tb_tcm":             s.HasilPemeriksaanTbTcm,
		"hasil_pemeriksaan_dm":                 s.HasilPemeriksaanDm,
		"hasil_pemeriksaan_ht":                 s.HasilPemeriksaanHt,
		"hasil_pemeriksaan_tb_poct":            s.HasilPemeriksaanPoct,
		"hasil_pemeriksaan_tb_radiologi":       s.HasilPemeriksaanRadiologi,
		"inveksi_hiv_aids":                     s.InfeksiHivAids,
		"gejala_dan_tanda_batuk":               s.GejalaDanTandaBatuk,
		"gejala_dan_tanda_bb_turun":            s.GejalaDanTandaBbTurun,
		"gejala_dan_tanda_demam_hilang_timbul": s.GejalaDanTandaDemamHilangTimbul,
		"gejala_dan_tanda_lesu_malaise":        s.GejalaDanTandaLesuMalaise,
		"gejala_dan_tanda_berkeringat_malam":   s.GejalaDanTandaBerkeringatMalam,
		"gejala_dan_tanda_pembesaran_getah_bening": s.GejalaDanTandaPembesaranKelenjarGB,
		"kontak_pasien_tbc":                        s.KontakPasienTbc,
		"gejala_dan_tanda_tbc":                     s.GejalaDanTandaTbc,
		"tindak_lanjut_penegakan_diagnosa":         s.TindakLanjutPenegakanDiagnosa,
		"updated_at":                               s.UpdatedAt,
	}
}

// ToSkriningCKGResult converts SkriningCKGRaw to SkriningCKGResult
func (s *SkriningCKGRaw) ToSkriningCKGResult() SkriningCKGResult {
	result := SkriningCKGResult{
		// Identitas Pasien
		PasienCKGID:        s.PasienCKGID,
		PasienNIK:          s.PasienNIK,
		PasienNama:         s.PasienNama,
		PasienJenisKelamin: s.PasienJenisKelamin,
		PasienTglLahir:     s.PasienTglLahir,
		PasienUsia:         s.PasienUsia,
		PasienPekerjaan:    s.PasienPekerjaan,
		PasienAlamat:       s.PasienAlamat,
		PasienNoHandphone:  s.PasienNoHandphone,

		// Data Kunjungan
		TglPemeriksaan: s.TglPemeriksaan,

		// Data Hasil Pemeriksaan
		BeratBadan:  s.BeratBadan,
		TinggiBadan: s.TinggiBadan,
		StatusImt:   s.StatusImt,
		HasilGds:    s.HasilGds,
		HasilGdp:    s.HasilGdp,
		HasilGdpp:   s.HasilGdpp,

		// Data Faktor Risiko
		KekuranganGizi: s.KekuranganGizi,
		Merokok:        s.Merokok,
		PerokokPasif:   s.PerokokPasif,
		LansiaDiatas65: s.LansiaDiatas65,
		IbuHamil:       s.IbuHamil,
		InfeksiHivAids: s.InfeksiHivAids,
		RiwayatDm:      s.HasilPemeriksaanDm,
		RiwayatHt:      s.HasilPemeriksaanHt,

		// Skrining gejala dan tanda
		GejalaBatuk:                s.GejalaDanTandaBatuk,
		GejalaBbTurun:              s.GejalaDanTandaBbTurun,
		GejalaDemamHilangTimbul:    s.GejalaDanTandaDemamHilangTimbul,
		GejalaLesuMalaise:          s.GejalaDanTandaLesuMalaise,
		GejalaBerkeringatMalam:     s.GejalaDanTandaBerkeringatMalam,
		GejalaPembesaranKelenjarGB: s.GejalaDanTandaPembesaranKelenjarGB,
		KontakPasienTbc:            s.KontakPasienTbc,

		// Pemeriksaan Lab TB
		HasilPemeriksaanTbBta:     s.HasilPemeriksaanTbBta,
		HasilPemeriksaanTbTcm:     s.HasilPemeriksaanTbTcm,
		HasilPemeriksaanPoct:      s.HasilPemeriksaanPoct,
		HasilPemeriksaanRadiologi: s.HasilPemeriksaanRadiologi,
	}

	return result
}

// Hasil Data Skrining yang dikirim ke client
type SkriningCKGResult struct {
	// Identitas Pasien
	PasienCKGID              string  `json:"pasien_ckg_id"`
	PasienNIK                string  `json:"pasien_nik"`
	PasienNama               string  `json:"pasien_nama"`
	PasienJenisKelamin       string  `json:"pasien_jenis_kelamin"`
	PasienTglLahir           string  `json:"pasien_tgl_lahir"`
	PasienUsia               int     `json:"pasien_usia"`
	PasienPekerjaan          *string `bson:"pasien_pekerjaan"` //TODO: saat ini data belum tersedia di dwh mongodb
	PasienProvinsiSatusehat  *string `json:"pasien_provinsi_satusehat"`
	PasienKabkotaSatusehat   *string `json:"pasien_kabkota_satusehat"`
	PasienKecamatanSatusehat *string `json:"pasien_kecamatan_satusehat"`
	PasienKelurahanSatusehat *string `json:"pasien_kelurahan_satusehat"`
	PasienProvinsiSitb       *string `json:"pasien_provinsi_sitb"`
	PasienKabkotaSitb        *string `json:"pasien_kabkota_sitb"`
	PasienKecamatanSitb      *string `json:"pasien_kecamatan_sitb"`
	PasienKelurahanSitb      *string `json:"pasien_kelurahan_sitb"`
	PasienAlamat             *string `json:"pasien_alamat"`
	PasienNoHandphone        string  `json:"pasien_no_handphone"`

	// Data Kunjungan
	KodeFaskesSatusehat *string `json:"periksa_faskes_satusehat"`
	KodeFaskesSITB      *string `json:"periksa_faskes_sitb"`
	TglPemeriksaan      string  `json:"periksa_tgl"`

	// Data Hasil Pemeriksaan
	BeratBadan  *float64 `json:"hasil_berat_badan"`
	TinggiBadan *float64 `json:"hasil_tinggi_badan"`
	StatusImt   *string  `json:"hasil_imt"`
	HasilGds    *float64 `json:"hasil_gds"`
	HasilGdp    *float64 `json:"hasil_gdp"`
	HasilGdpp   *float64 `json:"hasil_gdpp"`

	// Data Faktor Risiko
	KekuranganGizi *string `json:"risiko_kekurangan_gizi"`
	Merokok        *string `json:"risiko_merokok"`
	PerokokPasif   *string `json:"risiko_perokok_pasif"`
	LansiaDiatas65 *string `json:"risiko_lansia"`
	IbuHamil       *string `json:"risiko_ibu_hamil"`
	RiwayatDm      *string `json:"risiko_dm"`
	RiwayatHt      *string `json:"risiko_hipertensi"`
	InfeksiHivAids *string `json:"risiko_hiv_aids"`

	// Skrining gejala dan tanda
	GejalaBatuk                *string `json:"gejala_batuk"`
	GejalaBbTurun              *string `json:"gejala_bb_turun"`                // dipersiapkan untuk variabel baru/tambahan Andak & Dewasa
	GejalaDemamHilangTimbul    *string `json:"gejala_demam_hilang_timbul"`     // dipersiapkan untuk variabel baru/tambahan Andak & Dewasa
	GejalaLesuMalaise          *string `json:"gejala_lesu_malaise"`            // dipersiapkan untuk variabel baru/tambahan Anak
	GejalaBerkeringatMalam     *string `json:"gejala_berkeringat_malam"`       // dipersiapkan untuk variabel baru/tambahan Dewasa
	GejalaPembesaranKelenjarGB *string `json:"gejala_pembesaran_getah_bening"` // dipersiapkan untuk variabel baru/tambahan Dewasa
	KontakPasienTbc            *string `json:"kontak_pasien_tbc"`
	HasilSkriningTbc           *string `json:"hasil_skrining_tbc"`
	TerdugaTb                  *string `json:"terduga_tb"`

	// Pemeriksaan Lab TB
	HasilPemeriksaanTbBta     *string `json:"pemeriksaan_tb_bta"`
	HasilPemeriksaanTbTcm     *string `json:"pemeriksaan_tb_tcm"`
	HasilPemeriksaanPoct      *string `json:"pemeriksaan_tb_poct"`
	HasilPemeriksaanRadiologi *string `json:"pemeriksaan_tb_radiologi"`
}

type SkriningCKGOutput struct {
	Count     int                 `json:"totalRecords"`
	PageTotal int                 `json:"totalPage"`
	PageSize  int                 `json:"sizePerPage"`
	Page      int                 `json:"currentPage"`
	Results   []SkriningCKGResult `json:"results"`
}

// FromMap creates a SkriningCKGResult from a map
func (s *SkriningCKGResult) FromMap(data map[string]any) {
	// Map fields
	if val, ok := data["pasien_ckg_id"].(string); ok {
		s.PasienCKGID = val
	}
	if val, ok := data["pasien_nik"].(string); ok {
		s.PasienNIK = val
	}
	if val, ok := data["pasien_nama"].(string); ok {
		s.PasienNama = val
	}
	if val, ok := data["pasien_jenis_kelamin"].(string); ok {
		s.PasienJenisKelamin = val
	}
	if val, ok := data["pasien_tgl_lahir"].(string); ok {
		s.PasienTglLahir = val
	}
	if val, ok := data["pasien_usia"].(float64); ok {
		s.PasienUsia = int(val)
	}
	if val, ok := data["pasien_pekerjaan"].(string); ok {
		s.PasienPekerjaan = &val
	}
	if val, ok := data["pasien_provinsi_satusehat"].(string); ok {
		s.PasienProvinsiSatusehat = &val
	}
	if val, ok := data["pasien_kabkota_satusehat"].(string); ok {
		s.PasienKabkotaSatusehat = &val
	}
	if val, ok := data["pasien_kecamatan_satusehat"].(string); ok {
		s.PasienKecamatanSatusehat = &val
	}
	if val, ok := data["pasien_kelurahan_satusehat"].(string); ok {
		s.PasienKelurahanSatusehat = &val
	}
	if val, ok := data["pasien_provinsi_sitb"].(string); ok {
		s.PasienProvinsiSitb = &val
	}
	if val, ok := data["pasien_kabkota_sitb"].(string); ok {
		s.PasienKabkotaSitb = &val
	}
	if val, ok := data["pasien_kecamatan_sitb"].(string); ok {
		s.PasienKecamatanSitb = &val
	}
	if val, ok := data["pasien_kelurahan_sitb"].(string); ok {
		s.PasienKelurahanSitb = &val
	}
	if val, ok := data["pasien_alamat"].(string); ok {
		s.PasienAlamat = &val
	}
	if val, ok := data["pasien_no_handphone"].(string); ok {
		s.PasienNoHandphone = val
	}
	if val, ok := data["periksa_faskes_satusehat"].(string); ok {
		s.KodeFaskesSatusehat = &val
	}
	if val, ok := data["periksa_faskes_sitb"].(string); ok {
		s.KodeFaskesSITB = &val
	}
	if val, ok := data["periksa_tgl"].(string); ok {
		s.TglPemeriksaan = val
	}
	if val, ok := data["hasil_berat_badan"].(float64); ok {
		s.BeratBadan = &val
	}
	if val, ok := data["hasil_tinggi_badan"].(float64); ok {
		s.TinggiBadan = &val
	}
	if val, ok := data["hasil_imt"].(string); ok {
		s.StatusImt = &val
	}
	if val, ok := data["hasil_gds"].(float64); ok {
		s.HasilGds = &val
	}
	if val, ok := data["hasil_gdp"].(float64); ok {
		s.HasilGdp = &val
	}
	if val, ok := data["hasil_gdpp"].(float64); ok {
		s.HasilGdpp = &val
	}
	if val, ok := data["risiko_kekurangan_gizi"].(string); ok {
		s.KekuranganGizi = &val
	}
	if val, ok := data["risiko_merokok"].(string); ok {
		s.Merokok = &val
	}
	if val, ok := data["risiko_perokok_pasif"].(string); ok {
		s.PerokokPasif = &val
	}
	if val, ok := data["risiko_lansia"].(string); ok {
		s.LansiaDiatas65 = &val
	}
	if val, ok := data["risiko_ibu_hamil"].(string); ok {
		s.IbuHamil = &val
	}
	if val, ok := data["risiko_dm"].(string); ok {
		s.RiwayatDm = &val
	}
	if val, ok := data["risiko_hipertensi"].(string); ok {
		s.RiwayatHt = &val
	}
	if val, ok := data["risiko_hiv_aids"].(string); ok {
		s.InfeksiHivAids = &val
	}
	if val, ok := data["gejala_batuk"].(string); ok {
		s.GejalaBatuk = &val
	}
	if val, ok := data["gejala_bb_turun"].(string); ok {
		s.GejalaBbTurun = &val
	}
	if val, ok := data["gejala_demam_hilang_timbul"].(string); ok {
		s.GejalaDemamHilangTimbul = &val
	}
	if val, ok := data["gejala_lesu_malaise"].(string); ok {
		s.GejalaLesuMalaise = &val
	}
	if val, ok := data["gejala_berkeringat_malam"].(string); ok {
		s.GejalaBerkeringatMalam = &val
	}
	if val, ok := data["gejala_pembesaran_getah_bening"].(string); ok {
		s.GejalaPembesaranKelenjarGB = &val
	}
	if val, ok := data["kontak_pasien_tbc"].(string); ok {
		s.KontakPasienTbc = &val
	}
	if val, ok := data["hasil_skrining_tbc"].(string); ok {
		s.HasilSkriningTbc = &val
	}
	if val, ok := data["terduga_tb"].(string); ok {
		s.TerdugaTb = &val
	}
	if val, ok := data["pemeriksaan_tb_bta"].(string); ok {
		s.HasilPemeriksaanTbBta = &val
	}
	if val, ok := data["pemeriksaan_tb_tcm"].(string); ok {
		s.HasilPemeriksaanTbTcm = &val
	}
	if val, ok := data["pemeriksaan_tb_radiologi"].(string); ok {
		s.HasilPemeriksaanRadiologi = &val
	}
	if val, ok := data["pemeriksaan_tb_poct"].(string); ok {
		s.HasilPemeriksaanPoct = &val
	}
}

func (s *SkriningCKGResult) ToMap() map[string]any {
	return map[string]any{
		"pasien_ckg_id":                  s.PasienCKGID,
		"pasien_nik":                     s.PasienNIK,
		"pasien_nama":                    s.PasienNama,
		"pasien_jenis_kelamin":           s.PasienJenisKelamin,
		"pasien_tgl_lahir":               s.PasienTglLahir,
		"pasien_usia":                    s.PasienUsia,
		"pasien_pekerjaan":               s.PasienPekerjaan,
		"pasien_provinsi_satusehat":      s.PasienProvinsiSatusehat,
		"pasien_kabkota_satusehat":       s.PasienKabkotaSatusehat,
		"pasien_kecamatan_satusehat":     s.PasienKecamatanSatusehat,
		"pasien_kelurahan_satusehat":     s.PasienKelurahanSatusehat,
		"pasien_provinsi_sitb":           s.PasienProvinsiSitb,
		"pasien_kabkota_sitb":            s.PasienKabkotaSitb,
		"pasien_kecamatan_sitb":          s.PasienKecamatanSitb,
		"pasien_kelurahan_sitb":          s.PasienKelurahanSitb,
		"pasien_alamat":                  s.PasienAlamat,
		"pasien_no_handphone":            s.PasienNoHandphone,
		"periksa_faskes_satusehat":       s.KodeFaskesSatusehat,
		"periksa_faskes_sitb":            s.KodeFaskesSITB,
		"periksa_tgl":                    s.TglPemeriksaan,
		"hasil_berat_badan":              s.BeratBadan,
		"hasil_tinggi_badan":             s.TinggiBadan,
		"hasil_imt":                      s.StatusImt,
		"hasil_gds":                      s.HasilGds,
		"hasil_gdp":                      s.HasilGdp,
		"hasil_gdpp":                     s.HasilGdpp,
		"risiko_kekurangan_gizi":         s.KekuranganGizi,
		"risiko_merokok":                 s.Merokok,
		"risiko_perokok_pasif":           s.PerokokPasif,
		"risiko_lansia":                  s.LansiaDiatas65,
		"risiko_ibu_hamil":               s.IbuHamil,
		"risiko_dm":                      s.RiwayatDm,
		"risiko_hipertensi":              s.RiwayatHt,
		"risiko_hiv_aids":                s.InfeksiHivAids,
		"gejala_batuk":                   s.GejalaBatuk,
		"gejala_bb_turun":                s.GejalaBbTurun,
		"gejala_demam_hilang_timbul":     s.GejalaDemamHilangTimbul,
		"gejala_lesu_malaise":            s.GejalaLesuMalaise,
		"gejala_berkeringat_malam":       s.GejalaBerkeringatMalam,
		"gejala_pembesaran_getah_bening": s.GejalaPembesaranKelenjarGB,
		"kontak_pasien_tbc":              s.KontakPasienTbc,
		"hasil_skrining_tbc":             s.HasilSkriningTbc,
		"terduga_tb":                     s.TerdugaTb,
		"pemeriksaan_tb_bta":             s.HasilPemeriksaanTbBta,
		"pemeriksaan_tb_tcm":             s.HasilPemeriksaanTbTcm,
		"pemeriksaan_tb_poct":            s.HasilPemeriksaanPoct,
		"pemeriksaan_tb_radiologi":       s.HasilPemeriksaanRadiologi,
	}
}
