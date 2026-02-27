-- Удаление поля total_salary из department_applications
ALTER TABLE department_applications DROP COLUMN IF EXISTS total_salary;
