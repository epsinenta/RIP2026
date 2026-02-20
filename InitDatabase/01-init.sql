CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    login VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    is_moderator BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS departments (
    department_id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(1000) NOT NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    photo_url VARCHAR(255),
    employee_count INTEGER NOT NULL DEFAULT 0,
    head VARCHAR(255),
    reports_to VARCHAR(255),
    video VARCHAR(255),
    short_description VARCHAR(500)
);

CREATE TABLE IF NOT EXISTS department_applications (
    department_application_id SERIAL PRIMARY KEY,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    creator_id INTEGER NOT NULL REFERENCES users(user_id),
    forming_date TIMESTAMP,
    finish_date TIMESTAMP,
    moderator_id INTEGER REFERENCES users(user_id),
    title VARCHAR(255),
    total_salary NUMERIC(12,2)
);

CREATE TABLE IF NOT EXISTS department_application_departments (
    department_application_id INTEGER NOT NULL REFERENCES department_applications(department_application_id),
    department_id INTEGER NOT NULL REFERENCES departments(department_id),
    amount INTEGER NOT NULL DEFAULT 1,
    level INTEGER NOT NULL DEFAULT 1,
    is_main BOOLEAN DEFAULT FALSE,
    role VARCHAR(100),
    salary NUMERIC(12,2),
    PRIMARY KEY (department_application_id, department_id)
);

INSERT INTO users (login, password, is_moderator) VALUES ('user1', 'pass1', false), ('moderator', 'modpass', true) ON CONFLICT (login) DO NOTHING;

INSERT INTO departments (title, description, is_deleted, photo_url, employee_count, head, reports_to, video, short_description) VALUES
('Отдел информационных технологий', 'Отдел информационных технологий является ключевым структурным подразделением компании.', false, 'it_department.jpg', 42, 'Иванов Алексей Петрович', 'Генеральный директор', 'it_department.mp4', 'Поддержка корпоративных систем и инфраструктуры'),
('Бухгалтерия и финансовый контроль', 'Отдел бухгалтерии и финансового контроля обеспечивает ведение бухгалтерского и налогового учёта.', false, 'accounting.jpg', 18, 'Смирнова Елена Викторовна', 'Финансовый директор', 'accounting.mp4', 'Финансовый учёт и отчётность компании'),
('Отдел кадров', 'Отдел кадров осуществляет полный цикл управления персоналом.', false, 'hr_department.jpg', 12, 'Козлова Мария Дмитриевна', 'Директор по персоналу', 'hr_department.mp4', 'Управление персоналом, подбор и обучение сотрудников'),
('Юридический отдел', 'Юридический отдел обеспечивает правовую поддержку всех направлений деятельности компании.', false, 'legal.jpg', 9, 'Петров Дмитрий Сергеевич', 'Генеральный директор', 'legal.mp4', 'Правовая поддержка деятельности компании'),
('Отдел закупок', 'Отдел закупок занимается обеспечением компании всеми необходимыми материальными ресурсами.', false, 'procurement.jpg', 15, 'Новиков Андрей Владимирович', 'Коммерческий директор', 'procurement.mp4', 'Обеспечение компании материальными ресурсами'),
('Административно-хозяйственный отдел', 'Административно-хозяйственный отдел отвечает за эксплуатацию и содержание офисных помещений.', false, 'admin_office.jpg', 20, 'Фёдоров Игорь Николаевич', 'Заместитель генерального директора', 'admin_office.mp4', 'Эксплуатация офисов и внутренняя логистика');
