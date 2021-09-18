import {Component, OnInit} from '@angular/core';
import { Post } from "../model/post";
import {GQLService} from "../service/gql.service";
import {DataLibService} from "data-lib";

@Component({
  selector: 'app-post',
  templateUrl: './post.component.html',
  styleUrls: ['./post.component.scss']
})
export class PostComponent implements OnInit {
  posts: Post[] = [];

  constructor(private gqlService: GQLService,
              private dataService: DataLibService) {
  }

  ngOnInit(): void {
    if (this.dataService.getPosts()) {
      this.posts = this.dataService.getPosts();
    }

    this.dataService.getPostSubject().subscribe(posts => {
      this.posts = posts;
    });
  }

  getMorePosts() {
    this.gqlService.refreshPosts()
  }
}
