package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

// ========== MODELS ==========

type Level struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	NameAr    string    `json:"name_ar"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

type Year struct {
	ID        int       `json:"id"`
	LevelID   int       `json:"level_id"`
	Name      string    `json:"name"`
	NameAr    string    `json:"name_ar"`
	CreatedAt time.Time `json:"created_at"`
	LevelName string    `json:"level_name,omitempty"`
}

type Subject struct {
	ID        int       `json:"id"`
	YearID    int       `json:"year_id"`
	Name      string    `json:"name"`
	NameAr    string    `json:"name_ar"`
	Icon      string    `json:"icon"`
	CreatedAt time.Time `json:"created_at"`
	YearName  string    `json:"year_name,omitempty"`
}

type Category struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	NameAr    string    `json:"name_ar"`
	CreatedAt time.Time `json:"created_at"`
}

type Document struct {
	ID           int       `json:"id"`
	SubjectID    int       `json:"subject_id"`
	CategoryID   int       `json:"category_id"`
	Title        string    `json:"title"`
	FileName     string    `json:"file_name"`
	FilePath     string    `json:"file_path"`
	FileSize     int64     `json:"file_size"`
	Downloads    int       `json:"downloads"`
	CreatedAt    time.Time `json:"created_at"`
	SubjectName  string    `json:"subject_name,omitempty"`
	CategoryName string    `json:"category_name,omitempty"`
}

var db *sql.DB

// ========== DATABASE INIT ==========

func initDB() error {
	var err error
	db, err = sql.Open("sqlite", "./StudyDz.db")
	if err != nil {
		return err
	}

	// Create tables
	queries := []string{
		`CREATE TABLE IF NOT EXISTS levels (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            name_ar TEXT NOT NULL,
            color TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )`,
		`CREATE TABLE IF NOT EXISTS years (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            level_id INTEGER NOT NULL,
            name TEXT NOT NULL,
            name_ar TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (level_id) REFERENCES levels(id)
        )`,
		`CREATE TABLE IF NOT EXISTS subjects (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            year_id INTEGER NOT NULL,
            name TEXT NOT NULL,
            name_ar TEXT NOT NULL,
            icon TEXT,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (year_id) REFERENCES years(id)
        )`,
		`CREATE TABLE IF NOT EXISTS categories (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            name_ar TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )`,
		`CREATE TABLE IF NOT EXISTS documents (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            subject_id INTEGER NOT NULL,
            category_id INTEGER NOT NULL,
            title TEXT NOT NULL,
            file_name TEXT NOT NULL,
            file_path TEXT NOT NULL,
            file_size INTEGER DEFAULT 0,
            downloads INTEGER DEFAULT 0,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (subject_id) REFERENCES subjects(id),
            FOREIGN KEY (category_id) REFERENCES categories(id)
        )`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	log.Println("✅ Database initialized successfully")
	return insertDefaultData()
}

func insertDefaultData() error {
	// Check if data exists
	var count int
	db.QueryRow("SELECT COUNT(*) FROM levels").Scan(&count)
	if count > 0 {
		log.Println("ℹ️  Default data already exists")
		return nil
	}

	// Insert Levels
	levels := []struct {
		name, nameAr, color string
	}{
		{"Primaire", "التعليم الابتدائي", "#ef4444"},
		{"Moyen", "التعليم المتوسط", "#10b981"},
		{"Lycée", "التعليم الثانوي", "#f59e0b"},
		{"Université", "الجامعة", "#06b6d4"},
	}

	for _, l := range levels {
		_, err := db.Exec("INSERT INTO levels (name, name_ar, color) VALUES (?, ?, ?)", l.name, l.nameAr, l.color)
		if err != nil {
			return err
		}
	}

	// Insert Years for Primaire (5 years)
	primaireYears := []string{
		"السنة الأولى ابتدائي",
		"السنة الثانية ابتدائي",
		"السنة الثالثة ابتدائي",
		"السنة الرابعة ابتدائي",
		"السنة الخامسة ابتدائي",
	}
	for i, year := range primaireYears {
		_, err := db.Exec("INSERT INTO years (level_id, name, name_ar) VALUES (?, ?, ?)",
			1, fmt.Sprintf("Année %d primaire", i+1), year)
		if err != nil {
			return err
		}
	}

	// Insert Years for Moyen (4 years)
	moyenYears := []string{
		"السنة الأولى متوسط",
		"السنة الثانية متوسط",
		"السنة الثالثة متوسط",
		"السنة الرابعة متوسط",
	}
	for i, year := range moyenYears {
		_, err := db.Exec("INSERT INTO years (level_id, name, name_ar) VALUES (?, ?, ?)",
			2, fmt.Sprintf("Année %d moyen", i+1), year)
		if err != nil {
			return err
		}
	}

	// Insert Years for Lycée (3 years)
	lyceeYears := []string{
		"السنة الأولى ثانوي",
		"السنة الثانية ثانوي",
		"السنة الثالثة ثانوي",
	}
	for i, year := range lyceeYears {
		_, err := db.Exec("INSERT INTO years (level_id, name, name_ar) VALUES (?, ?, ?)",
			3, fmt.Sprintf("Année %d secondaire", i+1), year)
		if err != nil {
			return err
		}
	}

	// Insert Categories
	categories := []struct {
		name, nameAr string
	}{
		{"Cours", "دروس"},
		{"Examens", "اختبارات"},
		{"Exercices", "تمارين"},
		{"Compositions", "فروض"},
		{"Résumés", "ملخصات"},
	}

	for _, c := range categories {
		_, err := db.Exec("INSERT INTO categories (name, name_ar) VALUES (?, ?)", c.name, c.nameAr)
		if err != nil {
			return err
		}
	}

	// Insert Subjects for ALL 5 years of Primaire
	primaireSubjects := []struct {
		name   string
		nameAr string
		icon   string
	}{
		{"Mathématiques", "الرياضيات", "📐"},
		{"Arabe", "اللغة العربية", "📖"},
		{"Français", "اللغة الفرنسية", "🇫🇷"},
		{"Anglais", "اللغة الإنجليزية", "🇬🇧"},
		{"Éducation Islamique", "التربية الإسلامية", "✨"},
		{"Sciences et Technologie", "التربية العلمية والتكنولوجية", "🔬"},
		{"Arts", "التربية الفنية", "🎨"},
		{"Éducation Civique", "التربية المدنية", "🏛️"},
		{"Musique", "التربية الموسيقية", "🎵"},
		{"Amazigh", "اللغة الأمازيغية", "ⵣ"},
		{"Activités", "أنشطة", "🖥️"},
		{"Écriture", "تعلم الكتابة", "✏️"},
		{"Fichiers divers", "ملفات متنوعة", "📁"},
		{"Chaînes YouTube", "قنوات يوتيوب", "▶️"},
		{"Conseils", "نصائح وتوجيهات", "💡"},
		{"Page principale", "الصفحة الرئيسية", "🏠"},
	}

	// Add subjects for each year of Primaire (1 to 5)
	for yearID := 1; yearID <= 5; yearID++ {
		for _, s := range primaireSubjects {
			_, err := db.Exec("INSERT INTO subjects (year_id, name, name_ar, icon) VALUES (?, ?, ?, ?)",
				yearID, s.name, s.nameAr, s.icon)
			if err != nil {
				log.Printf("Warning: Could not insert subject %s for year %d: %v", s.nameAr, yearID, err)
			}
		}
	}

	// Insert Subjects for ALL 4 years of Moyen
	moyenSubjects := []struct {
		name   string
		nameAr string
		icon   string
	}{
		{"Mathématiques", "الرياضيات", "📐"},
		{"Arabe", "اللغة العربية", "📖"},
		{"Français", "اللغة الفرنسية", "🇫🇷"},
		{"Anglais", "اللغة الإنجليزية", "🇬🇧"},
		{"Éducation Islamique", "التربية الإسلامية", "✨"},
		{"Histoire et Géographie", "التاريخ والجغرافيا", "🌍"},
		{"Sciences de la Nature et de la Vie", "علوم الطبيعة والحياة", "🔬"},
		{"Sciences Physiques", "العلوم الفيزيائية", "⚗️"},
		{"Éducation Civique", "التربية المدنية", "🏛️"},
		{"Arts", "التربية الفنية", "🎨"},
		{"Amazigh", "اللغة الأمازيغية", "ⵣ"},
		{"Informatique", "الإعلام الآلي", "💻"},
		{"Musique", "التربية الموسيقية", "🎵"},
		{"Chaînes YouTube", "قنوات يوتيوب", "▶️"},
		{"Calculateur de moyenne", "برنامج حساب المعدل", "🧮"},
		{"Page principale", "الصفحة الرئيسية", "🏠"},
		{"Conseils", "نصائح وتوجيهات", "💡"},
	}

	// Add subjects for each year of Moyen (6 to 9)
	for yearID := 6; yearID <= 9; yearID++ {
		for _, s := range moyenSubjects {
			_, err := db.Exec("INSERT INTO subjects (year_id, name, name_ar, icon) VALUES (?, ?, ?, ?)",
				yearID, s.name, s.nameAr, s.icon)
			if err != nil {
				log.Printf("Warning: Could not insert subject %s for year %d: %v", s.nameAr, yearID, err)
			}
		}
	}

	// Insert Subjects for 1ère année Lycée (year_id = 10)
	lycee1Subjects := []struct {
		name   string
		nameAr string
		icon   string
	}{
		{"Mathématiques", "الرياضيات", "📐"},
		{"Arabe", "اللغة العربية", "📖"},
		{"Français", "اللغة الفرنسية", "🇫🇷"},
		{"Anglais", "اللغة الإنجليزية", "🇬🇧"},
		{"Éducation Islamique", "التربية الإسلامية", "✨"},
		{"Histoire et Géographie", "التاريخ والجغرافيا", "🌍"},
		{"Sciences de la Nature et de la Vie", "علوم الطبيعة والحياة", "🔬"},
		{"Sciences Physiques", "العلوم الفيزيائية", "⚗️"},
		{"Technologie", "التكنولوجيا", "⚙️"},
		{"Informatique", "الإعلام الآلي", "💻"},
		{"Amazigh", "اللغة الأمازيغية", "ⵣ"},
		{"Arts", "التربية الفنية", "🎨"},
		{"Fichiers divers", "ملفات متنوعة", "📁"},
		{"Chaînes YouTube", "قنوات يوتيوب", "▶️"},
		{"Calculateur de moyenne", "برنامج حساب المعدل", "🧮"},
		{"Page principale", "الصفحة الرئيسية", "🏠"},
		{"Conseils", "نصائح وتوجيهات", "💡"},
	}

	for _, s := range lycee1Subjects {
		_, err := db.Exec("INSERT INTO subjects (year_id, name, name_ar, icon) VALUES (?, ?, ?, ?)",
			10, s.name, s.nameAr, s.icon)
		if err != nil {
			log.Printf("Warning: Could not insert subject %s for year 10: %v", s.nameAr, err)
		}
	}

	// Insert Subjects for 2ème année Lycée (year_id = 11)
	lycee2Subjects := []struct {
		name   string
		nameAr string
		icon   string
	}{
		{"Mathématiques", "الرياضيات", "📐"},
		{"Arabe", "اللغة العربية", "📖"},
		{"Français", "اللغة الفرنسية", "🇫🇷"},
		{"Anglais", "اللغة الإنجليزية", "🇬🇧"},
		{"Éducation Islamique", "التربية الإسلامية", "✨"},
		{"Histoire et Géographie", "التاريخ والجغرافيا", "🌍"},
		{"Sciences de la Nature et de la Vie", "علوم الطبيعة والحياة", "🔬"},
		{"Sciences Physiques", "العلوم الفيزيائية", "⚗️"},
		{"Gestion Comptable et Financière", "التسيير المحاسبي والمالي", "📊"},
		{"Économie et Management", "الإقتصاد والمناجمنت", "📈"},
		{"Droit", "القانون", "⚖️"},
		{"Génie Civil", "الهندسة المدنية", "🏗️"},
		{"Génie des Procédés", "هندسة الطرائق", "🔧"},
		{"Génie Mécanique", "الهندسة الميكانيكية", "⚙️"},
		{"Génie Électrique", "الهندسة الكهربائية", "⚡"},
		{"Espagnol", "اللغة الإسبانية", "🇪🇸"},
		{"Allemand", "اللغة الألمانية", "🇩🇪"},
		{"Amazigh", "اللغة الأمازيغية", "ⵣ"},
		{"Italien", "اللغة الإيطالية", "🇮🇹"},
		{"Philosophie", "الفلسفة", "🤔"},
		{"Arts", "التربية الفنية", "🎨"},
		{"Fichiers divers", "ملفات متنوعة", "📁"},
		{"Chaînes YouTube", "قنوات يوتيوب", "▶️"},
		{"Calculateur de moyenne", "برنامج حساب المعدل", "🧮"},
		{"Page principale", "الصفحة الرئيسية", "🏠"},
		{"Conseils", "نصائح وتوجيهات", "💡"},
	}

	for _, s := range lycee2Subjects {
		_, err := db.Exec("INSERT INTO subjects (year_id, name, name_ar, icon) VALUES (?, ?, ?, ?)",
			11, s.name, s.nameAr, s.icon)
		if err != nil {
			log.Printf("Warning: Could not insert subject %s for year 11: %v", s.nameAr, err)
		}
	}

	// Insert Subjects for 3ème année Lycée (year_id = 12)
	lycee3Subjects := []struct {
		name   string
		nameAr string
		icon   string
	}{
		{"Mathématiques", "الرياضيات", "📐"},
		{"Arabe", "اللغة العربية", "📖"},
		{"Français", "اللغة الفرنسية", "🇫🇷"},
		{"Anglais", "اللغة الإنجليزية", "🇬🇧"},
		{"Éducation Islamique", "التربية الإسلامية", "✨"},
		{"Histoire et Géographie", "التاريخ والجغرافيا", "🌍"},
		{"Sciences de la Nature et de la Vie", "علوم الطبيعة والحياة", "🔬"},
		{"Sciences Physiques", "العلوم الفيزيائية", "⚗️"},
		{"Philosophie", "الفلسفة", "🤔"},
		{"Gestion Comptable et Financière", "التسيير المحاسبي والمالي", "📊"},
		{"Économie et Management", "الإقتصاد والمناجمنت", "📈"},
		{"Droit", "القانون", "⚖️"},
		{"Génie Civil", "الهندسة المدنية", "🏗️"},
		{"Génie des Procédés", "هندسة الطرائق", "🔧"},
		{"Génie Mécanique", "الهندسة الميكانيكية", "⚙️"},
		{"Génie Électrique", "الهندسة الكهربائية", "⚡"},
		{"Espagnol", "اللغة الإسبانية", "🇪🇸"},
		{"Allemand", "اللغة الألمانية", "🇩🇪"},
		{"Amazigh", "اللغة الأمازيغية", "ⵣ"},
		{"Italien", "اللغة الإيطالية", "🇮🇹"},
		{"Portail Universitaire", "بوابة التعليم الجامعي", "🎓"},
		{"Guide du Baccalauréat", "مواضيع ودليل شهادة الباكالوريا", "📝"},
		{"Chaînes YouTube", "قنوات يوتيوب", "▶️"},
		{"Calculateur de moyenne", "برنامج حساب المعدل", "🧮"},
		{"Page principale", "الصفحة الرئيسية", "🏠"},
		{"Conseils", "نصائح وتوجيهات", "💡"},
	}

	for _, s := range lycee3Subjects {
		_, err := db.Exec("INSERT INTO subjects (year_id, name, name_ar, icon) VALUES (?, ?, ?, ?)",
			12, s.name, s.nameAr, s.icon)
		if err != nil {
			log.Printf("Warning: Could not insert subject %s for year 12: %v", s.nameAr, err)
		}
	}

	log.Println("✅ Default data inserted")
	log.Println("📚 Primaire: 5 years x 16 subjects = 80 subjects")
	log.Println("📚 Moyen: 4 years x 17 subjects = 68 subjects")
	log.Println("📚 Lycée 1ère: 17 subjects")
	log.Println("📚 Lycée 2ème: 26 subjects")
	log.Println("📚 Lycée 3ème: 26 subjects")
	log.Println("📊 Total: 217 subjects")
	return nil
}

// ========== PUBLIC API HANDLERS ==========

func GetLevels(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, name_ar, color, created_at FROM levels ORDER BY id")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var levels []Level
	for rows.Next() {
		var l Level
		if err := rows.Scan(&l.ID, &l.Name, &l.NameAr, &l.Color, &l.CreatedAt); err != nil {
			continue
		}
		levels = append(levels, l)
	}
	c.JSON(200, levels)
}

func GetYears(c *gin.Context) {
	levelID := c.Query("level_id")

	query := `SELECT y.id, y.level_id, y.name, y.name_ar, y.created_at, l.name_ar as level_name 
              FROM years y 
              JOIN levels l ON y.level_id = l.id 
              WHERE y.level_id = ? 
              ORDER BY y.id`

	rows, err := db.Query(query, levelID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var years []Year
	for rows.Next() {
		var y Year
		if err := rows.Scan(&y.ID, &y.LevelID, &y.Name, &y.NameAr, &y.CreatedAt, &y.LevelName); err != nil {
			continue
		}
		years = append(years, y)
	}
	c.JSON(200, years)
}

func GetSubjects(c *gin.Context) {
	yearID := c.Query("year_id")

	query := `SELECT s.id, s.year_id, s.name, s.name_ar, s.icon, s.created_at, y.name_ar as year_name 
              FROM subjects s 
              JOIN years y ON s.year_id = y.id 
              WHERE s.year_id = ? 
              ORDER BY s.id`

	rows, err := db.Query(query, yearID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var subjects []Subject
	for rows.Next() {
		var s Subject
		if err := rows.Scan(&s.ID, &s.YearID, &s.Name, &s.NameAr, &s.Icon, &s.CreatedAt, &s.YearName); err != nil {
			continue
		}
		subjects = append(subjects, s)
	}
	c.JSON(200, subjects)
}

func GetCategories(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, name_ar, created_at FROM categories ORDER BY id")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var cat Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.NameAr, &cat.CreatedAt); err != nil {
			continue
		}
		categories = append(categories, cat)
	}
	c.JSON(200, categories)
}

func GetDocuments(c *gin.Context) {
	subjectID := c.Query("subject_id")

	query := `SELECT d.id, d.subject_id, d.category_id, d.title, d.file_name, d.file_path, 
              d.file_size, d.downloads, d.created_at, s.name_ar as subject_name, cat.name_ar as category_name
              FROM documents d
              JOIN subjects s ON d.subject_id = s.id
              JOIN categories cat ON d.category_id = cat.id
              WHERE d.subject_id = ?
              ORDER BY d.created_at DESC`

	rows, err := db.Query(query, subjectID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var documents []Document
	for rows.Next() {
		var doc Document
		if err := rows.Scan(&doc.ID, &doc.SubjectID, &doc.CategoryID, &doc.Title, &doc.FileName,
			&doc.FilePath, &doc.FileSize, &doc.Downloads, &doc.CreatedAt, &doc.SubjectName, &doc.CategoryName); err != nil {
			continue
		}
		documents = append(documents, doc)
	}
	c.JSON(200, documents)
}

func GetStats(c *gin.Context) {
	var stats struct {
		TotalLevels    int `json:"total_levels"`
		TotalYears     int `json:"total_years"`
		TotalSubjects  int `json:"total_subjects"`
		TotalDocuments int `json:"total_documents"`
		TotalDownloads int `json:"total_downloads"`
	}

	db.QueryRow("SELECT COUNT(*) FROM levels").Scan(&stats.TotalLevels)
	db.QueryRow("SELECT COUNT(*) FROM years").Scan(&stats.TotalYears)
	db.QueryRow("SELECT COUNT(*) FROM subjects").Scan(&stats.TotalSubjects)
	db.QueryRow("SELECT COUNT(*) FROM documents").Scan(&stats.TotalDocuments)
	db.QueryRow("SELECT COALESCE(SUM(downloads), 0) FROM documents").Scan(&stats.TotalDownloads)

	c.JSON(200, stats)
}

func DownloadDocument(c *gin.Context) {
	docID := c.Param("id")

	var doc Document
	err := db.QueryRow("SELECT id, file_path, file_name FROM documents WHERE id = ?", docID).
		Scan(&doc.ID, &doc.FilePath, &doc.FileName)

	if err != nil {
		c.JSON(404, gin.H{"error": "Document not found"})
		return
	}

	db.Exec("UPDATE documents SET downloads = downloads + 1 WHERE id = ?", docID)
	c.FileAttachment(doc.FilePath, doc.FileName)
}

// ========== ADMIN API HANDLERS ==========

func GetAllYears(c *gin.Context) {
	query := `SELECT y.id, y.level_id, y.name, y.name_ar, y.created_at, l.name_ar as level_name 
              FROM years y 
              JOIN levels l ON y.level_id = l.id 
              ORDER BY y.level_id, y.id`

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var years []Year
	for rows.Next() {
		var y Year
		if err := rows.Scan(&y.ID, &y.LevelID, &y.Name, &y.NameAr, &y.CreatedAt, &y.LevelName); err != nil {
			continue
		}
		years = append(years, y)
	}
	c.JSON(200, years)
}

func GetAllSubjects(c *gin.Context) {
	query := `SELECT s.id, s.year_id, s.name, s.name_ar, s.icon, s.created_at, y.name_ar as year_name 
              FROM subjects s 
              JOIN years y ON s.year_id = y.id 
              ORDER BY s.year_id, s.id`

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var subjects []Subject
	for rows.Next() {
		var s Subject
		if err := rows.Scan(&s.ID, &s.YearID, &s.Name, &s.NameAr, &s.Icon, &s.CreatedAt, &s.YearName); err != nil {
			continue
		}
		subjects = append(subjects, s)
	}
	c.JSON(200, subjects)
}

func GetAllDocuments(c *gin.Context) {
	query := `SELECT d.id, d.subject_id, d.category_id, d.title, d.file_name, d.file_path, 
              d.file_size, d.downloads, d.created_at, s.name_ar as subject_name, cat.name_ar as category_name
              FROM documents d
              JOIN subjects s ON d.subject_id = s.id
              JOIN categories cat ON d.category_id = cat.id
              ORDER BY d.created_at DESC`

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var documents []Document
	for rows.Next() {
		var doc Document
		if err := rows.Scan(&doc.ID, &doc.SubjectID, &doc.CategoryID, &doc.Title, &doc.FileName,
			&doc.FilePath, &doc.FileSize, &doc.Downloads, &doc.CreatedAt, &doc.SubjectName, &doc.CategoryName); err != nil {
			continue
		}
		documents = append(documents, doc)
	}
	c.JSON(200, documents)
}

func CreateLevel(c *gin.Context) {
	var level Level
	if err := c.BindJSON(&level); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO levels (name, name_ar, color) VALUES (?, ?, ?)",
		level.Name, level.NameAr, level.Color)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	level.ID = int(id)
	c.JSON(201, level)
}

func UpdateLevel(c *gin.Context) {
	id := c.Param("id")
	var level Level
	if err := c.BindJSON(&level); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE levels SET name = ?, name_ar = ?, color = ? WHERE id = ?",
		level.Name, level.NameAr, level.Color, id)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Level updated successfully"})
}

func DeleteLevel(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("DELETE FROM levels WHERE id = ?", id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Level deleted successfully"})
}

func CreateYear(c *gin.Context) {
	var year Year
	if err := c.BindJSON(&year); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO years (level_id, name, name_ar) VALUES (?, ?, ?)",
		year.LevelID, year.Name, year.NameAr)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	year.ID = int(id)
	c.JSON(201, year)
}

func UpdateYear(c *gin.Context) {
	id := c.Param("id")
	var year Year
	if err := c.BindJSON(&year); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE years SET level_id = ?, name = ?, name_ar = ? WHERE id = ?",
		year.LevelID, year.Name, year.NameAr, id)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Year updated successfully"})
}

func DeleteYear(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("DELETE FROM years WHERE id = ?", id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Year deleted successfully"})
}

func CreateSubject(c *gin.Context) {
	var subject Subject
	if err := c.BindJSON(&subject); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO subjects (year_id, name, name_ar, icon) VALUES (?, ?, ?, ?)",
		subject.YearID, subject.Name, subject.NameAr, subject.Icon)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	subject.ID = int(id)
	c.JSON(201, subject)
}

func UpdateSubject(c *gin.Context) {
	id := c.Param("id")
	var subject Subject
	if err := c.BindJSON(&subject); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE subjects SET year_id = ?, name = ?, name_ar = ?, icon = ? WHERE id = ?",
		subject.YearID, subject.Name, subject.NameAr, subject.Icon, id)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Subject updated successfully"})
}

func DeleteSubject(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("DELETE FROM subjects WHERE id = ?", id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Subject deleted successfully"})
}

func CreateCategory(c *gin.Context) {
	var category Category
	if err := c.BindJSON(&category); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO categories (name, name_ar) VALUES (?, ?)",
		category.Name, category.NameAr)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	category.ID = int(id)
	c.JSON(201, category)
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category Category
	if err := c.BindJSON(&category); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE categories SET name = ?, name_ar = ? WHERE id = ?",
		category.Name, category.NameAr, id)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Category updated successfully"})
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("DELETE FROM categories WHERE id = ?", id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Category deleted successfully"})
}

func UploadDocument(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "No file uploaded"})
		return
	}

	subjectID := c.PostForm("subject_id")
	categoryID := c.PostForm("category_id")
	title := c.PostForm("title")

	os.MkdirAll("./uploads", 0755)

	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
	filePath := filepath.Join("uploads", filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(500, gin.H{"error": "Failed to save file"})
		return
	}

	_, err = db.Exec(`INSERT INTO documents (subject_id, category_id, title, file_name, file_path, file_size) 
                      VALUES (?, ?, ?, ?, ?, ?)`,
		subjectID, categoryID, title, file.Filename, filePath, file.Size)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to save to database"})
		return
	}

	c.JSON(200, gin.H{"message": "File uploaded successfully", "filename": filename})
}

func DeleteDocument(c *gin.Context) {
	docID := c.Param("id")

	var filePath string
	db.QueryRow("SELECT file_path FROM documents WHERE id = ?", docID).Scan(&filePath)

	_, err := db.Exec("DELETE FROM documents WHERE id = ?", docID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if filePath != "" {
		os.Remove(filePath)
	}

	c.JSON(200, gin.H{"message": "Document deleted successfully"})
}

// ========== MAIN ==========

func main() {
	if err := initDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	r.Static("/uploads", "./uploads")
	r.StaticFile("/image1.jpg", "./image1.jpg")
	r.StaticFile("/", "./index.html")
	r.StaticFile("/index.html", "./index.html")
	r.StaticFile("/level.html", "./level.html")
	r.StaticFile("/matiere.html", "./matiere.html")
	r.StaticFile("/documents.html", "./documents.html")
	r.StaticFile("/admin.html", "./admin.html")

	api := r.Group("/api")
	{
		// Public routes
		api.GET("/levels", GetLevels)
		api.GET("/years", GetYears)
		api.GET("/subjects", GetSubjects)
		api.GET("/categories", GetCategories)
		api.GET("/documents", GetDocuments)
		api.GET("/download/:id", DownloadDocument)
		api.GET("/stats", GetStats)

		// Admin routes - Get All
		api.GET("/admin/years", GetAllYears)
		api.GET("/admin/subjects", GetAllSubjects)
		api.GET("/admin/documents", GetAllDocuments)

		// Admin routes - Levels
		api.POST("/admin/levels", CreateLevel)
		api.PUT("/admin/levels/:id", UpdateLevel)
		api.DELETE("/admin/levels/:id", DeleteLevel)

		// Admin routes - Years
		api.POST("/admin/years", CreateYear)
		api.PUT("/admin/years/:id", UpdateYear)
		api.DELETE("/admin/years/:id", DeleteYear)

		// Admin routes - Subjects
		api.POST("/admin/subjects", CreateSubject)
		api.PUT("/admin/subjects/:id", UpdateSubject)
		api.DELETE("/admin/subjects/:id", DeleteSubject)

		// Admin routes - Categories
		api.POST("/admin/categories", CreateCategory)
		api.PUT("/admin/categories/:id", UpdateCategory)
		api.DELETE("/admin/categories/:id", DeleteCategory)

		// Admin routes - Documents
		api.POST("/admin/upload", UploadDocument)
		api.DELETE("/admin/documents/:id", DeleteDocument)
	}

	log.Println("✅ Database initialized successfully")
	log.Println("🚀 StudyDz Server running on http://localhost:8080")
	log.Println("📊 Stats: http://localhost:8080/api/stats")
	log.Println("⚙️  Admin: http://localhost:8080/admin.html")

	// Get port from environment variable (for Railway)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)

}
