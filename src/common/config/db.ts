export const db = {
  host: process.env.POSTGRES_HOST || 'localhost',
  user: process.env.POSTGRES_USER || 'development',
  password: process.env.POSTGRES_PASSWORD || 'development',
  database: process.env.POSTGRES_DATABASE || 'micrach',
};
