-- Migration: is_main -> main_department_id в м-м, поля по теме nullable
-- НЕ меняем department_applications и departments (только м-м)

-- 0. Удалить main_department_id из department_applications если был добавлен ранее (ошибочно)
ALTER TABLE department_applications DROP COLUMN IF EXISTS main_department_id;

-- 1. Добавить main_department_id и sort_order в department_application_departments
ALTER TABLE department_application_departments
    ADD COLUMN IF NOT EXISTS main_department_id INTEGER REFERENCES departments(department_id);

ALTER TABLE department_application_departments
    ADD COLUMN IF NOT EXISTS sort_order INTEGER NOT NULL DEFAULT 0;

-- 2. Перенести данные: если is_main=true, main_department_id=department_id (сам себе)
DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='department_application_departments' AND column_name='is_main') THEN
    UPDATE department_application_departments SET main_department_id = department_id WHERE is_main = TRUE AND main_department_id IS NULL;
  END IF;
END $$;

-- 3. Удалить is_main из м-м
ALTER TABLE department_application_departments
    DROP COLUMN IF EXISTS is_main;

-- 4. Поля по теме сделать nullable (role, salary, amount)
ALTER TABLE department_application_departments
    ALTER COLUMN amount DROP NOT NULL;

-- title и total_salary в department_applications уже nullable в init-up
-- role и salary в department_application_departments уже nullable по умолчанию
