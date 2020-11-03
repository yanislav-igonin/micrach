import { NotFoundException } from '../../common/exceptions';
import { PostData, PostsRepository } from './repositories/posts.repository';
import { ThreadsRepository } from './repositories/threads.repository';

export class ThreadsService {
  threadsRepository: ThreadsRepository;
  postsRepository: PostsRepository;

  constructor(
    threadsRepository: ThreadsRepository,
    postsRepository: PostsRepository,
  ) {
    this.threadsRepository = threadsRepository;
    this.postsRepository = postsRepository;
  }

  getAll(page: number) {
    return this.threadsRepository.getAll(page);
  }

  async getOne(id: number) {
    const thread = await this.threadsRepository.getOne(id);
    if (thread === undefined) {
      throw new NotFoundException('Thread not found');
    }
    return thread;
  }

  async createOne(data: Omit<PostData, 'threadId'>) {
    const thread = await this.threadsRepository.createOne();
    const post = await this.postsRepository.createOne({
      ...data,
      threadId: thread.id,
    });
    thread.posts = [post];
    return thread;
  }

  async createPost(threadId: number, data: Omit<PostData, 'threadId'>) {
    const thread = await this.threadsRepository.getOne(threadId);
    if (thread === undefined) {
      throw new NotFoundException('Thread not found');
    }

    await this.postsRepository.createOne({
      ...data,
      threadId: thread.id,
    });

    await this.threadsRepository.updateThreadTime(thread.id);

    return this.postsRepository.getThreadPosts(threadId);
  }
}
