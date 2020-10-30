export const db = {
  host: process.env.POSTGRES_HOST || 'localhost',
  user: process.env.POSTGRES_USER || 'development',
  password: process.env.POSTGRES_PASSWORD || 'development',
  db: process.env.POSTGRES_DB || 'micrach',
};
