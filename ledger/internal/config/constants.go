package config

const DataDirPath = "ledger/.data"

const TinkGwAddr = "localhost:8080"

const UnassignedExpensesAccount = "2000-01-01 open Expenses:Unassigned\n\n"

const UnassignedIncomeAccount = "2000-01-01 open Income:Unassigned\n\n"

const DefaultOperatingCurrency = "option \"operating_currency\" \"EUR\"\n\n"

const DownloadBufferSize = 128 * 1024 // 64KiB, tweak this as desired
