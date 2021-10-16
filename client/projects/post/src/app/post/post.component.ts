import {Component, OnDestroy, OnInit} from '@angular/core';
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
    this.gqlService.getPostCreated();
  }

  ngOnDestroy(): void {
    GQLService._oldestPostReached = false;
  }

  get oldestPostReached(): boolean {
    // TODO when changing components and coming back to this one
    // check whether it is necessary to refetch newer posts than the newest this client has
    // ideas: send newest post id to server and ask whether new posts are there
    // subscription notification maybe?
    //
    return GQLService._oldestPostReached;
  }

  onScroll() {
    if (!this.oldestPostReached) {
      this.gqlService.refreshPosts();
    }
  }
}
