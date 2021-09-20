import {Component, HostListener, OnDestroy, OnInit} from '@angular/core';
import { Post } from "../model/post";
import {GQLService} from "../service/gql.service";
import {Observable} from "rxjs";

@Component({
  selector: 'app-post',
  templateUrl: './post.component.html',
  styleUrls: ['./post.component.scss']
})
export class PostComponent implements OnInit, OnDestroy {
  posts: Observable<Post[]> | undefined;

  constructor(private gqlService: GQLService) {
  }

  ngOnInit(): void {
    this.posts = this.gqlService.getPosts();
  }

  ngOnDestroy(): void {
    this.gqlService.oldestPostReached = false;
  }

  @HostListener("window:scroll", [])
  onScroll(): void {
    if (this.bottomReached() && !this.oldestPostReached) {
      this.gqlService.refreshPosts();
    }
  }

  bottomReached(): boolean {
    return (window.innerHeight + window.scrollY) >= document.body.offsetHeight;
  }

  get oldestPostReached(): boolean {
    return this.gqlService.oldestPostReached;
  }
}
