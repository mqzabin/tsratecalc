module github.com/mqzabin/tsratecalc/shopspring

go 1.23.2

replace github.com/mqzabin/tsratecalc/basecalc => ../basecalc

require (
	github.com/mqzabin/tsratecalc/basecalc v0.0.0-00010101000000-000000000000
	github.com/shopspring/decimal v1.4.0
)
