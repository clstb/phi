package parser

func REs() map[string]string {
	return map[string]string{
		"Date":        DateRE.String(),
		"Account":     AccountRE.String(),
		"Amount":      AmountRE.String(),
		"Posting":     PostingRE.String(),
		"Transaction": TransactionRE.String(),
	}
}
