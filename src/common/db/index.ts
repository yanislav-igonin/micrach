import { getConnectionManager } from 'typeorm';
import { Post } from '../../modules/threads/entities/post.entity';
import { Thread } from '../../modules/threads/entities/thread.entity';
import { config } from '../config';

const connectionManager = getConnectionManager();
const connection = connectionManager.create({
  type: 'postgres',
  host: config.db.host,
  username: config.db.user,
  password: config.db.password,
  database: config.db.database,
  entities: [Thread, Post],
  // logging: true,
  synchronize: true,
});

export { connection as db };
