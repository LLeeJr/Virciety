import { Injectable } from '@angular/core';
import {Post} from "../../../post/src/app/model/post";

@Injectable({
  providedIn: 'root'
})
export class DataLibService {
  private _posts: Post[] = [];

  constructor() { }

  get posts(): Post[] {
    return this._posts;
  }

  set posts(posts) {
    this._posts = posts;
  }

  removePost(id: string) {
    this._posts = this._posts.filter(post => post.id !== id);
  }

  addNewPost(newPost: Post): boolean {
    this._posts.reverse();
    this._posts.push(newPost);
    this._posts.reverse();
    return true;
  }

  getPost(postID: string): Post | undefined {
    return this._posts.find(post => post.id === postID);
  }
}
