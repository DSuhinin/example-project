package fixtures

// LoadCurrencyExchangeData loads test fixtures for TestExchangeCurrency_OK test.
// nolint:lll
func (f Fixtures) LoadCurrencyExchangeData() error {
	_, err := f.connection.Exec(`
		INSERT INTO public.keys (
			id,
			value,
			created_at,
			updated_at
		) VALUES (
			1,
			'e12a1983-046a-4f2c-b5a2-e27a6851ec4c',
			'2020-08-27 14:49:35.000000 +00:00',
			'2018-08-10 10:24:33.000000 +00:00'
		)`,
	)
	if err != nil {
		return err
	}

	return nil
}
