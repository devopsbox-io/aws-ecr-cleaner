package main

import "testing"

func TestGetDefaultKeepDays(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		env      map[string]string
		expected int
	}{
		"Env variable not set": {
			env:      map[string]string{},
			expected: 30,
		},
		"Env variable with invalid value": {
			env: map[string]string{
				"DEFAULT_KEEP_DAYS": "invalid",
			},
			expected: 30,
		},
		"Env variable with valid value": {
			env: map[string]string{
				"DEFAULT_KEEP_DAYS": "10",
			},
			expected: 10,
		},
	}

	for name, testCase := range tests {
		// capture range variables
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := getDefaultKeepDays(testLookupEnv(testCase.env))

			if result != testCase.expected {
				t.Errorf("Result %v different than expected %v", result, testCase.expected)
			}
		})
	}
}

func TestGetDryRun(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		env      map[string]string
		expected bool
	}{
		"Env variable not set": {
			env:      map[string]string{},
			expected: true,
		},
		"Env variable with invalid value": {
			env: map[string]string{
				"DRY_RUN": "invalid",
			},
			expected: true,
		},
		"Env variable with valid true value": {
			env: map[string]string{
				"DRY_RUN": "true",
			},
			expected: true,
		},
		"Env variable with valid false value": {
			env: map[string]string{
				"DRY_RUN": "false",
			},
			expected: false,
		},
	}

	for name, testCase := range tests {
		// capture range variables
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := getDryRun(testLookupEnv(testCase.env))

			if result != testCase.expected {
				t.Errorf("Result %v different than expected %v", result, testCase.expected)
			}
		})
	}
}

func TestIsLambda(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		env      map[string]string
		expected bool
	}{
		"Env variable not set": {
			env:      map[string]string{},
			expected: false,
		},
		"Env variable set": {
			env: map[string]string{
				"AWS_LAMBDA_FUNCTION_NAME": "lambda1",
			},
			expected: true,
		},
	}

	for name, testCase := range tests {
		// capture range variables
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := isLambda(testLookupEnv(testCase.env))

			if result != testCase.expected {
				t.Errorf("Result %v different than expected %v", result, testCase.expected)
			}
		})
	}
}

func testLookupEnv(env map[string]string) func(key string) (string, bool) {
	return func(key string) (string, bool) {
		result, exists := env[key]
		return result, exists
	}
}
