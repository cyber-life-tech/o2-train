package drift

const (
	// KolmogorovSmirnov test excels for subtle drifts in continuous data on small datasets (<10k samples) or research,
	// offering high sensitivity without bucketing and p-value testing,
	// but risks over-alerting on large data and fails on categorical
	KolmogorovSmirnov Metric = "ks"
	// PopulationStabilityIndex is best for categorical/mixed data on large datasets in production,
	// providing stable, low-false-positive detection with thresholds (e.g., 0.1-0.2),
	// but lacks sensitivity for small shifts and relies on bucketing
	PopulationStabilityIndex Metric = "psi"
)

// Metric of a drift detection operation
type Metric string
