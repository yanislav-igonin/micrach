import { Post } from './entities/post.entity';
import { Thread } from './entities/thread.entity';
import { ThreadsService } from './threads.service';
import { ThreadsRepository } from './repositories/threads.repository';
import { PostsRepository } from './repositories/posts.repository';
import { NotFoundException } from '../../common/exceptions';

const getPostMock = (threadId: number, id: number): Post => ({
  id,
  threadId,
  title: 'MOCK TITLE',
  text: 'MOCK TEXT',
  isSage: false,
  files: [],
  createdAt: (new Date()).toISOString(),
});
const getThreadMock = (id: number): Thread => ({
  id,
  posts: [getPostMock(id, 1), getPostMock(id, 2)],
  createdAt: (new Date()).toISOString(),
  updatedAt: (new Date()).toISOString(),
});

describe('ThreadsService', () => {
  let service: ThreadsService;
  let threadsRepository: ThreadsRepository;
  let postsRepository: PostsRepository;

  beforeEach(async () => {
    threadsRepository = new ThreadsRepository();
    postsRepository = new PostsRepository();
    service = new ThreadsService(threadsRepository, postsRepository);
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });

  describe('get all', () => {
    it('should return 0 threads', async () => {
      jest
        .spyOn(threadsRepository, 'getAll')
        .mockResolvedValueOnce([]);

      const threads = await service.getAll();
      expect(threads.length).toEqual(0);
    });

    it('should return 2 threads', async () => {
      jest
        .spyOn(threadsRepository, 'getAll')
        .mockResolvedValueOnce([getThreadMock(1), getThreadMock(2)]);

      const threads = await service.getAll();
      expect(threads.length).toEqual(2);
    });
  });

  describe('get one', () => {
    it('should return thread', async () => {
      jest
        .spyOn(threadsRepository, 'getOne')
        .mockResolvedValueOnce(getThreadMock(1));

      const thread = await service.getOne(1);
      expect(thread).toBeDefined();
    });

    it('should throw NotFoundException for non-existent thread', async () => {
      jest
        .spyOn(threadsRepository, 'getOne')
        .mockResolvedValueOnce(undefined);

      try {
        await service.getOne(1);
        expect(true).toBeFalsy();
      } catch (err) {
        expect(err).toBeInstanceOf(NotFoundException);
        expect(err.message).toEqual('Thread not found');
      }
    });
  });
});
