package calculator

import (
	"testing"
)

const (
	MCI2025   = 3932.0
	Tolerance = 0.5 // Allow 0.5 tenge difference due to floating point
)

func TestCalculateFromGross(t *testing.T) {
	calc := NewCalculator(MCI2025)

	tests := []struct {
		name           string
		gross          float64
		expectedNetMin float64 // Минимальное значение NET
		expectedNetMax float64 // Максимальное значение NET
		expectError    bool
	}{
		{
			name:           "Standard salary 500,000",
			gross:          500000,
			expectedNetMin: 401000, // >= 401000
			expectedNetMax: 402500, // <= 402500
			expectError:    false,
		},
		{
			name:           "Minimum salary (90,001)",
			gross:          90001,
			expectedNetMin: 0,
			expectedNetMax: 100000,
			expectError:    false,
		},
		{
			name:        "Too small salary (90,000)",
			gross:       90000,
			expectError: true,
		},
		{
			name:           "High salary 1,000,000",
			gross:          1000000,
			expectedNetMin: 797000,
			expectedNetMax: 798500,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := calc.CalculateFromGross(tt.gross)

			if (err != nil) != tt.expectError {
				t.Errorf("CalculateFromGross() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if err != nil {
				return // Expected error
			}

			// Verify structure
			if resp.GrossSalary != tt.gross {
				t.Errorf("GrossSalary = %f, want %f", resp.GrossSalary, tt.gross)
			}

			// Check NET is in expected range
			if resp.NetSalary < tt.expectedNetMin || resp.NetSalary > tt.expectedNetMax {
				t.Errorf("NetSalary = %f, expected range [%f, %f]",
					resp.NetSalary, tt.expectedNetMin, tt.expectedNetMax)
			}

			// Verify deductions don't exceed salary
			totalDeductions := resp.OPV + resp.VOSMS + resp.IPN
			if totalDeductions > tt.gross {
				t.Errorf("Total deductions %.2f exceed gross %.2f", totalDeductions, tt.gross)
			}

			// Verify NET = Gross - Deductions (approx)
			calculatedNet := tt.gross - resp.OPV - resp.VOSMS - resp.IPN
			if diff := calculatedNet - resp.NetSalary; diff > Tolerance || diff < -Tolerance {
				t.Errorf("NET calculation error: expected %.2f, got %.2f (diff: %.2f)",
					calculatedNet, resp.NetSalary, diff)
			}

			// Verify tax percentages
			expectedOPV := tt.gross * 0.10
			if diff := expectedOPV - resp.OPV; diff > Tolerance || diff < -Tolerance {
				t.Errorf("OPV calculation error: expected %.2f, got %.2f", expectedOPV, resp.OPV)
			}

			expectedVOSMS := tt.gross * 0.02
			if diff := expectedVOSMS - resp.VOSMS; diff > Tolerance || diff < -Tolerance {
				t.Errorf("VOSMS calculation error: expected %.2f, got %.2f", expectedVOSMS, resp.VOSMS)
			}
		})
	}
}

func TestCalculateFromNet(t *testing.T) {
	calc := NewCalculator(MCI2025)

	tests := []struct {
		name        string
		net         float64
		expectError bool
	}{
		{
			name:        "Standard net 350,000",
			net:         350000,
			expectError: false,
		},
		{
			name:        "Minimum net (90,001)",
			net:         90001,
			expectError: false,
		},
		{
			name:        "Too small net (90,000)",
			net:         90000,
			expectError: true,
		},
		{
			name:        "High net 700,000",
			net:         700000,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := calc.CalculateFromNet(tt.net)

			if (err != nil) != tt.expectError {
				t.Errorf("CalculateFromNet() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if err != nil {
				return // Expected error
			}

			// Verify that calculated NET matches input NET (within tolerance)
			if diff := resp.NetSalary - tt.net; diff > Tolerance || diff < -Tolerance {
				t.Errorf("NetSalary = %.2f, want %.2f (diff: %.2f)",
					resp.NetSalary, tt.net, diff)
			}

			// Verify that we can reverse-calculate and get same NET
			reverseResp, err := calc.CalculateFromGross(resp.GrossSalary)
			if err != nil {
				t.Errorf("failed to reverse-calculate: %v", err)
				return
			}

			if diff := reverseResp.NetSalary - tt.net; diff > Tolerance || diff < -Tolerance {
				t.Errorf("Reverse calculation NetSalary = %.2f, want %.2f (diff: %.2f)",
					reverseResp.NetSalary, tt.net, diff)
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	// Test that Gross -> Net -> Gross gives consistent results
	calc := NewCalculator(MCI2025)
	originalGross := 500000.0

	// Gross -> Net
	resp1, err := calc.CalculateFromGross(originalGross)
	if err != nil {
		t.Fatalf("CalculateFromGross failed: %v", err)
	}

	// Net -> Gross
	resp2, err := calc.CalculateFromNet(resp1.NetSalary)
	if err != nil {
		t.Fatalf("CalculateFromNet failed: %v", err)
	}

	// Check if we get back approximately the same gross
	if diff := resp2.GrossSalary - originalGross; diff > Tolerance || diff < -Tolerance {
		t.Errorf("Round-trip error: original %.2f, final %.2f (diff: %.2f)",
			originalGross, resp2.GrossSalary, diff)
	}
}

func TestEmployerContributions(t *testing.T) {
	calc := NewCalculator(MCI2025)

	gross := 500000.0
	resp, err := calc.CalculateFromGross(gross)
	if err != nil {
		t.Fatalf("CalculateFromGross failed: %v", err)
	}

	// Verify employer contributions are positive
	if resp.SO < 0 || resp.OOSMS < 0 || resp.SN < 0 {
		t.Errorf("Negative employer contributions: SO=%.2f, OOSMS=%.2f, SN=%.2f",
			resp.SO, resp.OOSMS, resp.SN)
	}

	// Verify employer total includes all contributions
	expectedTotal := gross + resp.SO + resp.OOSMS + resp.SN
	if diff := expectedTotal - resp.EmployerTotal; diff > Tolerance || diff < -Tolerance {
		t.Errorf("Employer total mismatch: expected %.2f, got %.2f",
			expectedTotal, resp.EmployerTotal)
	}
}
