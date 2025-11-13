package utils

import (
	"context"
	"errors"
	"pubsub-ckg-tb/internal/db/connection"
	"pubsub-ckg-tb/internal/models"
	"strings"
)

func IsNotEmptyString(str *string) bool {
	return str != nil && *str != ""
}

func IsNotEmptyInt(integer *int) bool {
	return integer != nil && *integer != 0
}

func FindMasterWilayah(id string, ctx context.Context, db connection.DatabaseConnection, collectionName string, useCache bool, cache *map[string]models.MasterWilayah) (*models.MasterWilayah, error) {
	if useCache {
		if val, ok := (*cache)[id]; ok {
			// fmt.Printf(" --------> Ambil Wilayah: %+v\n", val)
			return &val, nil
		}
	}

	var level int
	level = 1
	id = strings.ReplaceAll(id, ".", "")
	depdagriID := id
	ln := len(id)
	if ln == 10 {
		level = 4
		depdagriID = depdagriID[0:2] + "." + depdagriID[2:4] + "." + depdagriID[4:6] + "." + depdagriID[6:]
	} else if ln == 6 {
		level = 3
		depdagriID = depdagriID[0:2] + "." + depdagriID[2:4] + "." + depdagriID[4:6]
	} else if ln == 4 {
		level = 2
		depdagriID = depdagriID[0:2] + "." + depdagriID[2:4]
	}

	// fmt.Printf(" ----> cek wilayah %s (%s), len: %d, level: %d\n", depdagriID, id, ln, level)

	// filter := bson.D{
	// 	// {Key: "level", Value: level},
	// 	{Key: "id", Value: depdagriID},
	// }
	filter := map[string]any{
		// "level": level,
		"id": depdagriID,
	}

	var masterWilayah models.MasterWilayah
	err := db.FindOne(ctx, &masterWilayah, collectionName, nil, filter, nil)
	// err := collection.FindOne(ctx, filter).Decode(&masterWilayah)
	if err != nil {
		return nil, err
	}

	level = masterWilayah.Level

	if useCache {
		// fmt.Printf(" --------> Cached level wilayah %d\n", level)
		(*cache)[id] = masterWilayah
	}

	// fmt.Printf(" --------> Wilayah: %+v\n", masterWilayah)
	if level > 3 {
		FindMasterWilayah(id[0:6], ctx, db, collectionName, useCache, cache)
	}
	if level > 2 {
		FindMasterWilayah(id[0:4], ctx, db, collectionName, useCache, cache)
	}
	if level > 1 {
		FindMasterWilayah(id[0:2], ctx, db, collectionName, useCache, cache)
	}

	return &masterWilayah, nil
}

func FindMasterFaskes(id string, ctx context.Context, db connection.DatabaseConnection, collectionName string, useCache bool, cache *map[string]models.MasterFaskes) (*models.MasterFaskes, error) {
	if useCache {
		if val, ok := (*cache)[id]; ok {
			return &val, nil
		}
	}

	// fmt.Printf(" --------> Kode Faskes %s\n", id)

	// filter := bson.D{
	// 	{Key: "kode_satusehat", Value: id},
	// }
	filter := map[string]any{
		"kode_satusehat": id,
	}

	var masterFaskes models.MasterFaskes
	err := db.FindOne(ctx, &masterFaskes, collectionName, nil, filter, nil)
	// err := collection.FindOne(ctx, filter).Decode(&masterFaskes)
	if err != nil {
		// fmt.Printf(" --------> Error Faskes %s\n", err)
		return nil, err
	}
	// fmt.Printf(" --------> Faskes: %+v\n", masterFaskes)

	if useCache {
		(*cache)[id] = masterFaskes
	}

	return &masterFaskes, nil
}

func FindPasienTb(ctxStatusPasien context.Context, db connection.DatabaseConnection, collectionName string, item models.StatusPasien) (*models.StatusPasien, error) {
	orFilter := []map[string]any{}
	if IsNotEmptyString(item.PasienCkgID) {
		orFilter = append(orFilter, map[string]any{"pasien_ckg_id": item.PasienCkgID})
	} else {
		if IsNotEmptyString(item.TerdugaID) {
			orFilter = append(orFilter, map[string]any{"terduga_id": item.TerdugaID})
		}

		if IsNotEmptyString(item.PasienNIK) {
			orFilter = append(orFilter, map[string]any{"pasien_nik": item.PasienNIK})
		}
	}
	filter := map[string]any{
		"$or": orFilter,
	}

	var statusPasien models.StatusPasien
	err := db.FindOne(ctxStatusPasien, &statusPasien, collectionName, nil, filter, nil)
	// err := collectionStatusPasien.FindOne(ctxStatusPasien, filter).Decode(&statusPasien)
	if err != nil {
		return nil, err
	}

	return &statusPasien, nil
}

func InsertPasienTb(ctx context.Context, db connection.DatabaseConnection, collectionName string, item models.StatusPasien) (string, error) {
	// _, err := collection.InsertOne(ctx, item)
	_, err := db.InsertOne(ctx, collectionName, item)
	if err != nil {
		return err.Error(), err
	}

	return "new tb patient status added successfully", nil
}

func UpdatePasienTb(ctx context.Context, db connection.DatabaseConnection, collectionName string, item models.StatusPasien) (string, error) {
	setUpdate := map[string]any{
		"status_diagnosa":            item.StatusDiagnosis,
		"diagnosa_lab_hasil_tcm":     item.DiagnosisLabHasilTCM,
		"diagnosa_lab_hasil_bta":     item.DiagnosisLabHasilBTA,
		"tanggal_mulai_pengobatan":   item.TanggalMulaiPengobatan,
		"tanggal_selesai_pengobatan": item.TanggalSelesaiPengobatan,
		"hasil_akhir":                item.HasilAkhir,
	}
	orFilter := []map[string]any{}
	if IsNotEmptyString(item.PasienCkgID) {
		orFilter = append(orFilter, map[string]any{"pasien_ckg_id": item.PasienCkgID})
		setUpdate["pasien_ckg_id"] = item.PasienCkgID

		if IsNotEmptyString(item.TerdugaID) {
			setUpdate["terduga_id"] = item.TerdugaID
		}

		if IsNotEmptyString(item.PasienNIK) {
			setUpdate["pasien_nik"] = item.PasienNIK
		}
	} else {
		if IsNotEmptyString(item.TerdugaID) {
			orFilter = append(orFilter, map[string]any{"terduga_id": item.TerdugaID})
			setUpdate["terduga_id"] = item.TerdugaID
		}

		if IsNotEmptyString(item.PasienNIK) {
			orFilter = append(orFilter, map[string]any{"pasien_nik": item.PasienNIK})
			setUpdate["pasien_nik"] = item.PasienNIK
		}
	}

	if IsNotEmptyString(item.PasienTbID) {
		setUpdate["pasien_tb_id"] = item.PasienTbID
	}

	filter := map[string]any{
		"$or": orFilter,
	}

	// result, err := collection.UpdateOne(ctx, filter, update)
	result, err := db.UpdateOne(ctx, collectionName, filter, setUpdate)
	if err != nil {
		return err.Error(), err
	}

	if result == 0 {
		err = errors.New("failed to update tb patient status")
		return err.Error(), err
	}

	return "tb patient status updated successfully", nil
}
