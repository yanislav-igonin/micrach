import {
  Entity,
  PrimaryGeneratedColumn,
  CreateDateColumn,
  ManyToOne,
} from 'typeorm';
import { Thread } from './thread.entity';

@Entity('posts')
export class Post {
  @PrimaryGeneratedColumn()
  id!: number;

  @ManyToOne(() => Thread, (thread) => thread.posts)
  thread!: Thread;

  @CreateDateColumn()
  createdAt!: string;
}
