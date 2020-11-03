import { EntityRepository, Repository } from 'typeorm';
import { Thread } from '../entities/thread.entity';

@EntityRepository(Thread)
export class ThreadsRepository extends Repository<Thread> {
  getAll(page: number) {
    return this.find({
      skip: (page - 1) * 10,
      take: 10,
      relations: ['posts', 'posts.files'],
      order: { updatedAt: 'DESC' },
    });
  }

  getOne(id: number) {
    return this.findOne(id, { relations: ['posts', 'posts.files'] });
  }

  createOne() {
    return this.save({});
  }

  updateThreadTime(threadId: number) {
    return this.query(
      'UPDATE threads SET "updatedAt" = DEFAULT WHERE "id" IN ($1)',
      [threadId],
    );
  }
}
