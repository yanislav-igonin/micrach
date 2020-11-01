import { EntityRepository, Repository } from 'typeorm';
import { Post } from '../entities/post.entity';

export type PostData = Pick<Post, 'title' | 'text' | 'threadId' | 'isSage'>;

@EntityRepository(Post)
export class PostsRepository extends Repository<Post> {
  createOne(data: PostData) {
    return this.save(data);
  }
}
