package data

type Language struct {
	Name        string   `json:"name"`
	Paradigm    string   `json:"paradigm"`
	Compilation string   `json:"compilation"`
	Performance string   `json:"performance"`
	Concurrency string   `json:"concurrency"`
	ErrorModel  string   `json:"error_model"`
	UseCases    []string `json:"use_cases"`
	Companies   []string `json:"notable_companies_indonesia"`
	PopRank     int      `json:"popularity_rank_indonesia_2024"`
	LearnScore  int      `json:"ease_of_learning_score"`
	SpeedScore  int      `json:"speed_score"`
	EcoScore    int      `json:"ecosystem_score"`
}

var Languages = []Language{
	{
		Name: "Golang", Paradigm: "Procedural + minimal OOP",
		Compilation: "Native binary", Performance: "Very Fast",
		Concurrency: "Goroutines + Channels", ErrorModel: "Explicit error return values",
		UseCases:  []string{"High-performance backend API", "Microservices", "Cloud-native apps", "DevOps tools (Docker, K8s)"},
		Companies: []string{"Gojek", "Tokopedia", "Grab", "Halodoc", "Kumparan", "Ruangguru"},
		PopRank: 1, LearnScore: 8, SpeedScore: 10, EcoScore: 8,
	},
	{
		Name: "Python", Paradigm: "Multi-paradigm",
		Compilation: "Interpreted", Performance: "Moderate",
		Concurrency: "Threads / AsyncIO", ErrorModel: "Try-Except exceptions",
		UseCases:  []string{"Data Science & AI/ML", "Scripting & automation", "Web (Django, Flask)", "Research"},
		Companies: []string{"Traveloka", "Bukalapak", "Dana", "Tiket.com"},
		PopRank: 2, LearnScore: 10, SpeedScore: 5, EcoScore: 10,
	},
	{
		Name: "Java", Paradigm: "Classical OOP",
		Compilation: "JVM (bytecode)", Performance: "High",
		Concurrency: "Threads + Executor", ErrorModel: "Try-Catch exceptions",
		UseCases:  []string{"Enterprise applications", "Android development", "Banking systems", "Large-scale backend"},
		Companies: []string{"Bank BRI", "Mandiri", "BCA Tech", "Shopee Indonesia"},
		PopRank: 3, LearnScore: 6, SpeedScore: 8, EcoScore: 9,
	},
	{
		Name: "JavaScript", Paradigm: "Multi-paradigm",
		Compilation: "JIT (V8)", Performance: "Moderate-High",
		Concurrency: "Event Loop / Promise", ErrorModel: "Try-Catch / Promise.catch",
		UseCases:  []string{"Frontend web", "Backend (Node.js)", "Mobile (React Native)", "Full-stack"},
		Companies: []string{"Tokopedia", "Blibli", "OVO", "Gojek Frontend"},
		PopRank: 4, LearnScore: 9, SpeedScore: 6, EcoScore: 10,
	},
	{
		Name: "PHP", Paradigm: "Multi-paradigm",
		Compilation: "Interpreted", Performance: "Moderate",
		Concurrency: "Limited (FPM)", ErrorModel: "Try-Catch exceptions",
		UseCases:  []string{"Web development", "CMS (WordPress, Laravel)", "E-commerce backend"},
		Companies: []string{"Bukalapak (legacy)", "Kaskus", "Detik.com"},
		PopRank: 5, LearnScore: 9, SpeedScore: 5, EcoScore: 8,
	},
}
