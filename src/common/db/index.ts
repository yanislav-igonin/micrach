import { getConnectionManager } from 'typeorm';
import { config } from '../config';

const connectionManager = getConnectionManager();
const connection = connectionManager.create({
  type: 'postgres',
  host: config.db.host,
  username: config.db.user,
  password: config.db.password,
  database: config.db.database,
  // entities: [File, User],
  // logging: app.dbLog,
  // synchronize: app.dbSync,
});

export { connection as db };
