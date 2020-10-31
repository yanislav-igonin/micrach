import { EntityRepository, Repository } from 'typeorm';
import { Thread } from '../entities/thread.entity';

@EntityRepository(Thread)
export class ThreadsRepository extends Repository<Thread> {
  getAll() {
    return this.find({ relations: ['posts'] });
  }

  createOne() {
    return this.save({});
  }
}
