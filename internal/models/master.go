package models

type MasterFaskes struct {
	ID                   *string `bson:"id"`
	Nama                 *string `bson:"nama"`
	JenisUnitID          *string `bson:"jenis_unit_id"`
	KodeYankes           *string `bson:"kode_yankes"`
	KodeSatusehat        *string `bson:"kode_satusehat"`
	KodeSatusehatSandbox *string `bson:"kode_satusehat_sandbox"`
	JenisFasyankesID     *string `bson:"jenis_fasyankes_id"`
	TipeFasyankesID      *string `bson:"tipe_fasyankes_id"`
	JenisPemilikID       *string `bson:"jenis_pemilik_id"`
	StatusDotsID         *string `bson:"status_dots_id"`
	PelaksanaMikID       *string `bson:"pelaksana_mik_id"`
	ProvinsiID           *string `bson:"provinsi_id"`
	KabupatenID          *string `bson:"kabupaten_id"`
	KecamatanID          *string `bson:"kecamatan_id"`
	KelurahanID          *string `bson:"kelurahan_id"`
	Alamat               *string `bson:"alamat"`
	LocX                 *string `bson:"loc_x"`
	LocY                 *string `bson:"loc_y"`

	// Telepon                  *string `bson:"telepon"`
	// Fax                      *string `bson:"fax"`
	// KodePos                  *string `bson:"kode_pos"`
	// NamaKontak               *string `bson:"nama_kontak"`
	// Email                    *string `bson:"email"`
	// Website                  *string `bson:"website"`
	// StatusLabID              *string `bson:"status_lab_id"`
	// LabMikID                 *string `bson:"lab_mik_id"`
	// LabTcmID                 *string `bson:"lab_tcm_id"`
	// LabXdrID                 *string `bson:"lab_xdr_id"`
	// LabLpaLini1ID            *string `bson:"lab_lpa_lini1_id"`
	// LabLpaLini2ID            *string `bson:"lab_lpa_lini2_id"`
	// LabBiakanID              *string `bson:"lab_biakan_id"`
	// LabDstID                 *string `bson:"lab_dst_id"`

	StatusPerawatanID        *string `bson:"status_perawatan_id"`
	StatusGenexpertID        *string `bson:"status_genexpert_id"`
	StatusPengobatanID       *string `bson:"status_pengobatan_id"`
	StatusRujukanRoID        *string `bson:"status_rujukan_ro_id"`
	StatusSubRujukanRoID     *string `bson:"status_sub_rujukan_ro_id"`
	StatusPenyimpanID        *string `bson:"status_penyimpan_id"`
	StatusPenerimaObatID     *string `bson:"status_penerima_obat_id"`
	StatusPerubahanJmlObatID *string `bson:"status_perubahan_jml_obat_id"`
	StatusMintaLebihID       *string `bson:"status_minta_lebih_id"`
	StatusPemasokID          *string `bson:"status_pemasok_id"`
	PenyediaObatLini1ProvID  *string `bson:"penyedia_obat_lini1_prov_id"`
	PenyediaObatLini1ID      *string `bson:"penyedia_obat_lini1_id"`
	PenyediaObatLini2ProvID  *string `bson:"penyedia_obat_lini2_prov_id"`
	PenyediaObatLini2ID      *string `bson:"penyedia_obat_lini2_id"`

	// KodePmdt                 *string `bson:"kode_pmdt"`
	// StatusOnlineID           *string `bson:"status_online_id"`

	StatusAktifID *string `bson:"status_aktif_id"`

	// Keterangan               *string `bson:"keterangan"`
	// NoRekening               *string `bson:"no_rekening"`
	// PemilikRekening          *string `bson:"pemilik_rekening"`
	// NamaBank                 *string `bson:"nama_bank"`
	// CabangBank               *string `bson:"cabang_bank"`
	// TujuanKlaimLabID         *string `bson:"tujuan_klaim_lab_id"`
	// TujuanKlaimSuntikID      *string `bson:"tujuan_klaim_suntik_id"`
	// TujuanKlaimEnablerID     *string `bson:"tujuan_klaim_enabler_id"`
}

type MasterWilayah struct {
	ID            string  `bson:"id"`
	Kode          string  `bson:"kode"`
	Nama          string  `bson:"nama"`
	Level         int     `bson:"level"`
	ProvinsiID    *string `bson:"provinsi_id"`
	KodeProvinsi  *string `bson:"kode_provinsi"`
	KabupatenID   *string `bson:"kabupaten_id"`
	KodeKabupaten *string `bson:"kode_kabupaten"`
	KecamatanID   *string `bson:"kecamatan_id"`
	KodeKecamatan *string `bson:"kode_kecamatan"`
	KelurahanID   *string `bson:"kelurahan_id"`
	KodeKelurahan *string `bson:"kode_kelurahan"`
}
