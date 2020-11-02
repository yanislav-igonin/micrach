import {
  Entity,
  PrimaryGeneratedColumn,
  CreateDateColumn,
  ManyToOne,
  Column,
} from 'typeorm';
import { Post } from './post.entity';

@Entity('files')
export class File {
  @PrimaryGeneratedColumn()
  id!: number;

  @Column()
  postId!: number;

  @ManyToOne(() => Post, (post) => post.files)
  post?: Post;

  @Column()
  name!: string;

  @Column()
  extension!: string;

  @CreateDateColumn()
  createdAt!: string;
}
