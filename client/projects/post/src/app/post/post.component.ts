import {Component, OnInit} from '@angular/core';
import { Post } from "../model/post";
import {GQLService} from "../service/gql.service";

@Component({
  selector: 'app-post',
  templateUrl: './post.component.html',
  styleUrls: ['./post.component.scss']
})
export class PostComponent implements OnInit {
  loading = true;
  error: any;
  posts: Post[] = [];

  constructor(private gqlService: GQLService) {
  }

  ngOnInit(): void {
    this.gqlService.getPosts().subscribe((data: any) => {

      console.log(data.data.getPosts);

      for (let getPost of data.data.getPosts) {
        const post: Post = new Post(getPost);
        this.posts.push(post)
      }

      this.loading = data.loading;
      this.error = data.error;
    });
  }
}
