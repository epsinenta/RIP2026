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

type DepartmentApplication struct {
	ID             int
	Title          string
	Description    string
	AppCount       int
	TotalEmployees int
	TotalSalary    float64
	Hierarchy      []DepartmentApplicationDepartment
}

type DepartmentApplicationDepartment struct {
	Department Department
	Level      int
	Role       string
	Salary     float64
}

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

func (r *Repository) GetDepartmentApplications() ([]DepartmentApplication, error) {
	departments, err := r.GetDepartments()
	if err != nil {
		return nil, err
	}
	depMap := make(map[int]Department)
	for _, d := range departments {
		depMap[d.ID] = d
	}

	departmentApplications := []DepartmentApplication{
		{
			ID:             1,
			Title:          "Административная структура компании",
			Description:    "Формирование иерархической структуры подчинения отделов компании с определением уровней управления и расчётом итоговой зарплаты руководителей подразделений на основе количества подчинённых сотрудников.",
			AppCount:       6,
			TotalEmployees: 116,
			TotalSalary:    1380000,
			Hierarchy: []DepartmentApplicationDepartment{
				{Department: depMap[1], Level: 1, Role: "Головное подразделение", Salary: 410000},
				{Department: depMap[2], Level: 1, Role: "Руководящий отдел", Salary: 240000},
				{Department: depMap[3], Level: 2, Role: "Подчинённый отдел", Salary: 160000},
				{Department: depMap[4], Level: 2, Role: "Подчинённый отдел", Salary: 145000},
				{Department: depMap[5], Level: 3, Role: "Подчинённый отдел", Salary: 175000},
				{Department: depMap[6], Level: 2, Role: "Руководящий отдел", Salary: 250000},
			},
		},
	}
	return departmentApplications, nil
}

func (r *Repository) GetDepartmentApplication(id int) (DepartmentApplication, error) {
	apps, err := r.GetDepartmentApplications()
	if err != nil {
		return DepartmentApplication{}, err
	}

	for _, app := range apps {
		if app.ID == id {
			return app, nil
		}
	}
	return DepartmentApplication{}, fmt.Errorf("Заявка не найдена")
}

func (r *Repository) GetDepartmentApplicationForDepartment(departmentID int) (*DepartmentApplicationDepartment, error) {
	apps, err := r.GetDepartmentApplications()
	if err != nil {
		return nil, err
	}

	for _, app := range apps {
		for _, ad := range app.Hierarchy {
			if ad.Department.ID == departmentID {
				return &ad, nil
			}
		}
	}
	return nil, fmt.Errorf("Отдел не найден в заявке")
}
