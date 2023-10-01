package transaction

type CampaignTransactionFormatter struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Amount    int    `json:"amount"`
	CreatedAt string `json:"created_at"`
}

func FormatCampaignTransaction(transaction Transaction) CampaignTransactionFormatter {
	formatter := CampaignTransactionFormatter{}
	formatter.ID = transaction.ID
	formatter.Name = transaction.User.Name
	formatter.Amount = transaction.Amount
	formatter.CreatedAt = transaction.CreatedAt.Format("2006-01-02 15:04:05")

	return formatter
}

func FormatCampaignTransactions(transactions []Transaction) []CampaignTransactionFormatter {
	if len(transactions) == 0 {
		return []CampaignTransactionFormatter{}
	}

	var transactionsFormatter []CampaignTransactionFormatter

	for _, transaction := range transactions {
		formatter := FormatCampaignTransaction(transaction)
		transactionsFormatter = append(transactionsFormatter, formatter)
	}

	return transactionsFormatter
}
