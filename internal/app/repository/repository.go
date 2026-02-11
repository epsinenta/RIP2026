package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

// Department — услуга (отдел компании)
type Department struct {
	ID               int
	Title            string
	EmployeeCount    int
	Head             string
	ReportsTo        string
	Photo            string
	Video            string
	ShortDescription string
	Description      string
}

// Application — словарь заявки (административная структура)
type Application struct {
	ID             int
	Title          string
	Description    string
	Departments    []ApplicationDepartment
	AppCount       int
	TotalEmployees int
	TotalSalary    float64
}

// ApplicationDepartment — отдел в составе заявки (м-м связь)
type ApplicationDepartment struct {
	Department Department
	Level      int
	Role       string
	NewSalary  float64
}

// GetDepartments возвращает коллекцию всех отделов (услуг)
func (r *Repository) GetDepartments() ([]Department, error) {
	departments := []Department{
		{
			ID:               1,
			Title:            "Отдел информационных технологий",
			EmployeeCount:    42,
			Head:             "Иванов Алексей Петрович",
			ReportsTo:        "Генеральный директор",
			Photo:            "it_department.jpg",
			Video:            "it_department.mp4",
			ShortDescription: "Поддержка корпоративных систем и инфраструктуры",
			Description:      "Отдел информационных технологий является ключевым структурным подразделением компании, обеспечивающим бесперебойную работу всех корпоративных информационных систем. В зону ответственности отдела входит администрирование серверной инфраструктуры, поддержка локальной вычислительной сети, обеспечение информационной безопасности и защита данных. Отдел отвечает за внедрение и сопровождение программного обеспечения, используемого всеми подразделениями компании: ERP-систем, CRM-платформ, систем электронного документооборота и корпоративных порталов.",
		},
		{
			ID:               2,
			Title:            "Бухгалтерия и финансовый контроль",
			EmployeeCount:    18,
			Head:             "Смирнова Елена Викторовна",
			ReportsTo:        "Финансовый директор",
			Photo:            "accounting.jpg",
			Video:            "accounting.mp4",
			ShortDescription: "Финансовый учёт и отчётность компании",
			Description:      "Отдел бухгалтерии и финансового контроля обеспечивает ведение бухгалтерского и налогового учёта предприятия в полном соответствии с действующим законодательством. Сотрудники отдела осуществляют расчёт заработной платы, формирование финансовой отчётности, контроль движения денежных средств и взаимодействие с налоговыми органами. Отдел также отвечает за внутренний аудит и контроль за соблюдением финансовой дисциплины во всех подразделениях компании.",
		},
		{
			ID:               3,
			Title:            "Отдел кадров",
			EmployeeCount:    12,
			Head:             "Козлова Мария Дмитриевна",
			ReportsTo:        "Директор по персоналу",
			Photo:            "hr_department.jpg",
			Video:            "hr_department.mp4",
			ShortDescription: "Управление персоналом, подбор и обучение сотрудников",
			Description:      "Отдел кадров осуществляет полный цикл управления персоналом: от подбора и найма сотрудников до их адаптации, обучения и развития. Специалисты отдела разрабатывают кадровую политику, ведут кадровый учёт и делопроизводство, организуют программы повышения квалификации и корпоративные мероприятия. Отдел также занимается формированием корпоративной культуры и системы мотивации персонала.",
		},
		{
			ID:               4,
			Title:            "Юридический отдел",
			EmployeeCount:    9,
			Head:             "Петров Дмитрий Сергеевич",
			ReportsTo:        "Генеральный директор",
			Photo:            "legal.jpg",
			Video:            "legal.mp4",
			ShortDescription: "Правовая поддержка деятельности компании",
			Description:      "Юридический отдел обеспечивает правовую поддержку всех направлений деятельности компании. Сотрудники отдела осуществляют экспертизу и согласование договоров, представляют интересы компании в судебных инстанциях и государственных органах, консультируют руководство и сотрудников по правовым вопросам. Отдел разрабатывает внутренние нормативные документы и следит за соответствием деятельности компании действующему законодательству.",
		},
		{
			ID:               5,
			Title:            "Отдел закупок",
			EmployeeCount:    15,
			Head:             "Новиков Андрей Владимирович",
			ReportsTo:        "Коммерческий директор",
			Photo:            "procurement.jpg",
			Video:            "procurement.mp4",
			ShortDescription: "Обеспечение компании материальными ресурсами",
			Description:      "Отдел закупок занимается обеспечением компании всеми необходимыми материальными ресурсами, оборудованием и услугами сторонних поставщиков. Сотрудники отдела проводят анализ рынка, организуют тендерные процедуры, ведут переговоры с поставщиками и контролируют исполнение контрактов. Отдел стремится к оптимизации закупочных процессов и снижению затрат при сохранении высокого качества поставляемых товаров и услуг.",
		},
		{
			ID:               6,
			Title:            "Административно-хозяйственный отдел",
			EmployeeCount:    20,
			Head:             "Фёдоров Игорь Николаевич",
			ReportsTo:        "Заместитель генерального директора",
			Photo:            "admin_office.jpg",
			Video:            "admin_office.mp4",
			ShortDescription: "Эксплуатация офисов и внутренняя логистика",
			Description:      "Административно-хозяйственный отдел отвечает за эксплуатацию и содержание офисных помещений, организацию транспортного обслуживания, обеспечение охраны труда и пожарной безопасности. Сотрудники отдела координируют работу клининговых и охранных служб, контролируют техническое состояние инженерных систем зданий, организуют переезды и ремонтные работы. Отдел обеспечивает комфортные условия труда для всех подразделений компании.",
		},
	}
	if len(departments) == 0 {
		return nil, fmt.Errorf("Массив пустой")
	}

	return departments, nil
}

// GetDepartment возвращает отдел по ID
func (r *Repository) GetDepartment(id int) (Department, error) {
	departments, err := r.GetDepartments()
	if err != nil {
		return Department{}, err
	}

	for _, dep := range departments {
		if dep.ID == id {
			return dep, nil
		}
	}
	return Department{}, fmt.Errorf("Отдел не найден")
}

// GetDepartmentByTitle — фильтрация отделов по наименованию (поиск)
func (r *Repository) GetDepartmentByTitle(title string) ([]Department, error) {
	departments, err := r.GetDepartments()
	if err != nil {
		return []Department{}, err
	}

	var result []Department
	for _, dep := range departments {
		if strings.Contains(strings.ToLower(dep.Title), strings.ToLower(title)) {
			result = append(result, dep)
		}
	}
	return result, nil
}

// CalculateNewSalary — расчёт новой зарплаты руководителя подразделения
// Формула: Зарплата = БазовыйОклад + (Кол-воСотрудников * K)
// БазовыйОклад: подчинённый отдел — 100000, руководящий — 150000, головное — 200000
// K — коэффициент нагрузки — 5000
func CalculateNewSalary(employeeCount int, role string) float64 {
	K := 5000.0
	var baseSalary float64

	switch role {
	case "Головное подразделение":
		baseSalary = 200000
	case "Руководящий отдел":
		baseSalary = 150000
	default: // Подчинённый отдел
		baseSalary = 100000
	}

	return baseSalary + float64(employeeCount)*K
}

// buildApplication собирает одну заявку по её данным
func (r *Repository) buildApplication(id int, title string, description string, entries []struct {
	DepartmentID int
	Level        int
	Role         string
}) (Application, error) {
	departments, err := r.GetDepartments()
	if err != nil {
		return Application{}, err
	}

	depMap := make(map[int]Department)
	for _, d := range departments {
		depMap[d.ID] = d
	}

	var appDeps []ApplicationDepartment
	totalEmployees := 0
	totalSalary := 0.0
	for _, e := range entries {
		dep, ok := depMap[e.DepartmentID]
		if !ok {
			continue
		}
		salary := CalculateNewSalary(dep.EmployeeCount, e.Role)
		appDeps = append(appDeps, ApplicationDepartment{
			Department: dep,
			Level:      e.Level,
			Role:       e.Role,
			NewSalary:  salary,
		})
		totalEmployees += dep.EmployeeCount
		totalSalary += salary
	}

	return Application{
		ID:             id,
		Title:          title,
		Description:    description,
		Departments:    appDeps,
		AppCount:       len(appDeps),
		TotalEmployees: totalEmployees,
		TotalSalary:    totalSalary,
	}, nil
}

// GetApplications возвращает массив всех заявок (административных структур)
func (r *Repository) GetApplications() ([]Application, error) {
	entries1 := []struct {
		DepartmentID int
		Level        int
		Role         string
	}{
		{1, 1, "Головное подразделение"},
		{2, 1, "Руководящий отдел"},
		{3, 2, "Подчинённый отдел"},
		{4, 2, "Подчинённый отдел"},
		{5, 3, "Подчинённый отдел"},
		{6, 2, "Руководящий отдел"},
	}

	app1, err := r.buildApplication(
		1,
		"Административная структура компании",
		"Формирование иерархической структуры подчинения отделов компании с определением уровней управления и расчётом итоговой зарплаты руководителей подразделений на основе количества подчинённых сотрудников.",
		entries1,
	)
	if err != nil {
		return nil, err
	}

	return []Application{app1}, nil
}

// GetApplication возвращает заявку по ID
func (r *Repository) GetApplication(id int) (Application, error) {
	apps, err := r.GetApplications()
	if err != nil {
		return Application{}, err
	}

	for _, app := range apps {
		if app.ID == id {
			return app, nil
		}
	}
	return Application{}, fmt.Errorf("Заявка не найдена")
}

// GetApplicationForDepartment — найти информацию о заявке для конкретного отдела
func (r *Repository) GetApplicationForDepartment(departmentID int) (*ApplicationDepartment, error) {
	apps, err := r.GetApplications()
	if err != nil {
		return nil, err
	}

	for _, app := range apps {
		for _, ad := range app.Departments {
			if ad.Department.ID == departmentID {
				return &ad, nil
			}
		}
	}
	return nil, fmt.Errorf("Отдел не найден в заявке")
}
