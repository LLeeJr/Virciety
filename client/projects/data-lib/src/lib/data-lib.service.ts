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

  addNewPost(newPost: Post): boolean {
    // Check if post already exists (for createPost-Subscription)
    for (let post of this._posts) {
      if (post.id === newPost.id) {
        return false;
      }
    }

    this._posts.reverse();
    this._posts.push(newPost);
    this._posts.reverse();
    return true;
  }
}
