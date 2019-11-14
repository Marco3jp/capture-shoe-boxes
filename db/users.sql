CREATE USER `capture_shoe_boxes`@'localhost' IDENTIFIED BY 'example_password';
CREATE USER `diff_shoe_boxes`@'localhost' IDENTIFIED BY 'example_password';
CREATE USER `is_exist_researcher_api`@'localhost' IDENTIFIED BY 'example_password';

GRANT INSERT ON is_exist_researcher.capture TO `capture_shoe_boxes`@'localhost';
GRANT SELECT ON is_exist_researcher.capture TO `diff_shoe_boxes`@'localhost';
GRANT INSERT ON is_exist_researcher.shoe_box TO `diff_shoe_boxes`@'localhost';
GRANT SELECT ON is_exist_researcher.shoe_box TO `is_exist_researcher_api`@'localhost';
