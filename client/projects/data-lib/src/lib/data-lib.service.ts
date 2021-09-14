import { Injectable } from '@angular/core';
import {Subject} from "rxjs";

@Injectable({
  providedIn: 'root'
})
export class DataLibService {

  private posts: any;
  private _posts = new Subject<any>();

  constructor() { }

  getPosts(): any {
    return this.posts;
  }

  getPostSubject(): Subject<any> {
    return this._posts;
  }

  setPosts(posts: any) {
    this.posts = posts;
    this._posts.next(this.posts);
  }

  addPost(post: any)  {
    this.posts.reverse();
    this.posts.push(post);
    this.posts.reverse();
    this._posts.next(this.posts);
  }
}
