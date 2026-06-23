package uom

import (
	"errors"
	"soom-be-go/internal/domain"
)

var (
	// ErrIncompatibleUnit dikembalikan kalau dua unit tidak punya base_unit yang sama
	// (misal mencoba convert Gram ke Liter — tidak masuk akal).
	ErrIncompatibleUnit = errors.New("unit tidak bisa dikonversi: base unit berbeda")

	// ErrMissingConversionData dikembalikan kalau unit tidak punya data konversi lengkap.
	ErrMissingConversionData = errors.New("data konversi unit tidak lengkap")
)

// Convert mengkonversi quantity dari unit asal (from) ke unit tujuan (to).
// Formula: quantity_tujuan = quantity_asal * (factor_asal / factor_tujuan)
//
// Contoh: Convert(500, Gram, Kilogram) -> 0.5
// karena Gram.ConversionFactor = 0.001 dan Kilogram.ConversionFactor = 1.
func Convert(quantity float64, from domain.Uom, to domain.Uom) (float64, error) {
	// Kalau unit asal dan tujuan sama persis, tidak perlu konversi apa-apa.
	if from.Id == to.Id {
		return quantity, nil
	}

	if from.BaseUnit == nil || from.ConversionFactor == nil {
		return 0, ErrMissingConversionData
	}
	if to.BaseUnit == nil || to.ConversionFactor == nil {
		return 0, ErrMissingConversionData
	}

	// Dua unit hanya bisa dikonversi kalau base_unit-nya sama
	// (misal sama-sama berbasis "kg", atau sama-sama berbasis "L").
	if *from.BaseUnit != *to.BaseUnit {
		return 0, ErrIncompatibleUnit
	}

	// Hindari pembagian dengan nol kalau ada data conversion_factor yang salah/0.
	if *to.ConversionFactor == 0 {
		return 0, ErrMissingConversionData
	}

	result := quantity * (*from.ConversionFactor / *to.ConversionFactor)
	return result, nil
}
