import { EntityRepository, Repository } from 'typeorm';
import { Thread } from './thread.entity';

@EntityRepository(Thread)
export class ThreadsRepository extends Repository<Thread> {
  getAll() {
    return this.find();
  }

  createOne() {
    return this.save({});
  }
}
