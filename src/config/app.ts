export const app = {
  env: process.env.NODE_ENV || 'production',
  port: process.env.PORT !== '' && process.env.PORT !== undefined
    ? parseInt(process.env.PORT, 10)
    : 80,
};
