package models

import (
	"time"
)

// StatusPasien represents the patient status format
type StatusPasien struct {
	// PubSubObject
	PasienCkgID *string `json:"pasien_ckg_id" bson:"pasien_ckg_id"`
	TerdugaID   *string `json:"terduga_id" bson:"terduga_id"`     // ID Terduga TB
	PasienTbID  *string `json:"pasien_tb_id" bson:"pasien_tb_id"` // ID Pasien TB SO/RO jika sudah terkonfirmasi positif atau dirawat
	PasienNIK   *string `json:"pasien_nik" bson:"pasien_nik"`

	// Parameter dikendalikan dari SITB
	StatusDiagnosis          *string `json:"status_diagnosa" bson:"status_diagnosa"`               // ["TBC SO", "TBC RO", "Bukan TBC"]
	DiagnosisLabHasilTCM     *string `json:"diagnosa_lab_hasil_tcm" bson:"diagnosa_lab_hasil_tcm"` // TCM: ["not_detected", "rif_sen", "rif_res", "rif_indet", "invalid", "error", "no_result", "tdl"]
	DiagnosisLabHasilBTA     *string `json:"diagnosa_lab_hasil_bta" bson:"diagnosa_lab_hasil_bta"` // BTA: ["negatif", "positif"]
	TanggalMulaiPengobatan   *string `json:"tanggal_mulai_pengobatan" bson:"tanggal_mulai_pengobatan"`
	TanggalSelesaiPengobatan *string `json:"tanggal_selesai_pengobatan" bson:"tanggal_selesai_pengobatan"`
	HasilAkhir               *string `json:"hasil_akhir" bson:"hasil_akhir"` // ["Sembuh", "Pengobatan Lengkap", "Pengobatan Gagal", "Meninggal", "Putus berobat (lost to follow up)", "Tidak dievaluasi/pindah", "Gagal karena Perubahan Diagnosis"]
}

type StatusPasienResult struct {
	PasienCkgID *string `json:"pasien_ckg_id"` // ID CKG
	TerdugaID   *string `json:"terduga_id"`    // ID Terduga TB
	PasienTbID  *string `json:"pasien_tb_id"`  // ID Pasien TB SO/RO jika sudah dirawat
	PasienNIK   *string `json:"pasien_nik"`    // NIK
	IsError     bool    `json:"error"`         // menandakan apakah error atau tidak
	Respons     string  `json:"message"`       // pesan respon pemrosesan
}

// NewStatusPasien creates a new StatusPasien instance
func NewStatusPasien() *StatusPasien {
	return &StatusPasien{}
}

// FromMap creates a StatusPasien from a map
func (s *StatusPasien) FromMap(data map[string]interface{}) {
	// Map fields
	if val, ok := data["terduga_id"].(string); ok {
		s.TerdugaID = &val
	}
	if val, ok := data["pasien_nik"].(string); ok {
		s.PasienNIK = &val
	}
	if val, ok := data["pasien_tb_id"].(string); ok {
		s.PasienTbID = &val
	}
	if val, ok := data["status_diagnosa"].(string); ok {
		s.StatusDiagnosis = &val
	}
	if val, ok := data["diagnosa_lab_hasil_tcm"].(string); ok {
		s.DiagnosisLabHasilTCM = &val
	}
	if val, ok := data["diagnosa_lab_hasil_bta"].(string); ok {
		s.DiagnosisLabHasilBTA = &val
	}
	if val, ok := data["tanggal_mulai_pengobatan"].(string); ok {
		if _, err := time.Parse("2006-01-02", val); err == nil {
			s.TanggalMulaiPengobatan = &val
		}
	}
	if val, ok := data["tanggal_selesai_pengobatan"].(string); ok {
		if _, err := time.Parse("2006-01-02", val); err == nil {
			s.TanggalSelesaiPengobatan = &val
		}
	}
	if val, ok := data["hasil_akhir"].(string); ok {
		s.HasilAkhir = &val
	}
}

func (s *StatusPasien) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"pasien_ckg_id":              s.PasienCkgID,
		"terduga_id":                 s.TerdugaID,
		"pasien_nik":                 s.PasienNIK,
		"pasien_tb_id":               s.PasienTbID,
		"status_diagnosa":            s.StatusDiagnosis,
		"diagnosa_lab_hasil_tcm":     s.DiagnosisLabHasilTCM,
		"diagnosa_lab_hasil_bta":     s.DiagnosisLabHasilBTA,
		"tanggal_mulai_pengobatan":   s.TanggalMulaiPengobatan,
		"tanggal_selesai_pengobatan": s.TanggalSelesaiPengobatan,
		"hasil_akhir":                s.HasilAkhir,
	}

	return result
}
