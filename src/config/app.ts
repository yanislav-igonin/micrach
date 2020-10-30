export const app = {
  port: process.env.PORT !== '' && process.env.PORT !== undefined
    ? parseInt(process.env.PORT, 10)
    : 3000,
};
