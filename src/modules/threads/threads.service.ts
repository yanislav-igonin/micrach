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

  getAll() {
    return this.threadsRepository.getAll();
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
}
