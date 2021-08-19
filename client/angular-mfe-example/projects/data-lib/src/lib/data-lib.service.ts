import { Injectable } from '@angular/core';
import {Post} from "../model/post";
import {Comment} from "../model/comment";

@Injectable({
  providedIn: 'root'
})
export class DataLibService {
  private _posts: Post[] = [];
  private _comments: Map<string, Comment[]> = new Map<string, Comment[]>();

  constructor() {}

  get posts(): Post[] {
    return this._posts;
  }

  set posts(posts: Post[]) {
    this._posts = posts;
  }

  get comments(): Map<string, Comment[]> {
    return this._comments;
  }

  set comments(comments: Map<string, Comment[]>) {
    this._comments = comments;
  }

  getPost(postId: string): Post | null {
    for (let post of this._posts) {
      if (post.id === postId) {
        if (this._comments.has(postId)) {
          post.comments = this._comments.get(postId);
          return post;
        }
      }
    }

    return null;
  }

}
