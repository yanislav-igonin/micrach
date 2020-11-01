import {
  Entity,
  PrimaryGeneratedColumn,
  CreateDateColumn,
  ManyToOne,
  Column,
} from 'typeorm';
import { Thread } from './thread.entity';

@Entity('posts')
export class Post {
  @PrimaryGeneratedColumn()
  id!: number;

  @Column()
  threadId!: number;

  @ManyToOne(() => Thread, (thread) => thread.posts)
  thread?: Thread;

  @Column()
  title!: string;

  @Column()
  text!: string;

  @Column()
  isSage!: boolean;

  @CreateDateColumn()
  createdAt!: string;
}
