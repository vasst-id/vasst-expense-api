INSERT INTO "vasst_expense".currency (currency_code, currency_name, currency_symbol, currency_decimal_places, currency_status)
VALUES
  ('IDR', 'Indonesian Rupiah', 'Rp', 0, 1),
  ('USD', 'US Dollar', '$', 2, 1),
  ('EUR', 'Euro', '€', 2, 1),
  ('SGD', 'Singapore Dollar', 'S$', 2, 1),
  ('JPY', 'Japanese Yen', '¥', 0, 1);

INSERT INTO "vasst_expense".banks (bank_id, bank_name, bank_code, bank_logo_url, status)
VALUES
  (1, 'Bank Central Asia', 'BCA', 'https://logo.clearbit.com/bca.co.id', 1),
  (2, 'Bank Mandiri', 'MANDIRI', 'https://logo.clearbit.com/bankmandiri.co.id', 1),
  (3, 'Bank Negara Indonesia', 'BNI', 'https://logo.clearbit.com/bni.co.id', 1),
  (4, 'Bank Rakyat Indonesia', 'BRI', 'https://logo.clearbit.com/bri.co.id', 1),
  (5, 'Bank Syariah Indonesia', 'BSI', 'https://logo.clearbit.com/bsi.co.id', 1);

INSERT INTO "vasst_expense".taxonomy (label, value, type, type_label, status)
VALUES
  ('Debit', '1', 'account_type', 'Account Type', 1),
  ('Credit', '2', 'account_type', 'Account Type', 1),
  ('Savings', '3', 'account_type', 'Account Type', 1),
  ('Cash', '4', 'account_type', 'Account Type', 1),
  ('Shared', '5', 'account_type', 'Account Type', 1);

INSERT INTO "vasst_expense".taxonomy (label, value, type, type_label, status)
VALUES
  ('Income', '1', 'transaction_type', 'Transaction Type', 1),
  ('Expense', '2', 'transaction_type', 'Transaction Type', 1);

INSERT INTO "vasst_expense".taxonomy (label, value, type, type_label, status)
VALUES
  ('Debit/QRIS', '1', 'payment_method', 'Payment Method', 1),
  ('Credit', '2', 'payment_method', 'Payment Method', 1),
  ('Cash', '3', 'payment_method', 'Payment Method', 1),
  ('Transfer', '4', 'payment_method', 'Payment Method', 1);

INSERT INTO "vasst_expense".taxonomy (label, value, type, type_label, status)
VALUES
  ('Weekly', '1', 'period_type', 'Period Type', 1),
  ('Monthly', '2', 'period_type', 'Period Type', 1),
  ('Yearly', '3', 'period_type', 'Period Type', 1),
  ('Event', '4', 'period_type', 'Period Type', 1);

INSERT INTO "vasst_expense".taxonomy (label, value, type, type_label, status)
VALUES
  ('Paid', '1', 'credit_status', 'Credit Status', 1),
  ('Unpaid', '2', 'credit_status', 'Credit Status', 1);

INSERT INTO "vasst_expense".categories (name, description, icon, is_system_category)
VALUES
  ('Food & Beverage', 'Meals, groceries, snacks, and drinks', 'restaurant', true),
  ('Transportation', 'Public transport, fuel, ride-hailing, parking', 'directions_car', true),
  ('Utilities', 'Electricity, water, internet, phone', 'bolt', true),
  ('Shopping', 'Clothes, electronics, online shopping', 'shopping_cart', true),
  ('Health', 'Doctor, pharmacy, insurance', 'local_hospital', true),
  ('Entertainment', 'Movies, streaming, games, events', 'movie', true),
  ('Education', 'Tuition, books, courses', 'school', true),
  ('Travel', 'Flights, hotels, tours', 'flight', true),
  ('Other', 'Miscellaneous expenses', 'category', true);

INSERT INTO "vasst_expense".subscription_plan
  (subscription_plan_name, subscription_plan_description, subscription_plan_features, subscription_plan_price, subscription_plan_currency_id, subscription_plan_status)
VALUES
  (
    'Free',
    'Basic plan for personal use',
    '{\"max_workspaces\":1,\"max_accounts\":3,\"max_transactions\":1000}',
    0.00,
    1, -- Assuming 1 is IDR
    1
  ),
  (
    'Pro',
    'Advanced plan for power users',
    '{\"max_workspaces\":10,\"max_accounts\":20,\"max_transactions\":100000,\"priority_support\":true}',
    50000.00,
    1, -- Assuming 1 is IDR
    1
  );

  