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

  addComment(post: Post, newComment: string) {
    const addCommentRequest = {
      postID: post.id,
      comment: newComment,
      createdBy: 'user3', //this.authService.username
    };

    this.gqlService.addComment(post, addCommentRequest);
  }

  triggerEvent(post: Post, event: string) {
    if (event === 'like') {
      this.likePost(post);
    } else if (event === 'edit') {
      this.editPost(post);
    } else if (event === 'remove') {
      this.removePost(post);
    } else {
      this.showComments(post);
    }
  }

  likePost(post: Post) {
    // TODO uncomment this
    const username: string = /*this.authService.userName*/ 'user4';
    let liked: boolean = true;

    // check if it's a like or unlike
    const index = post.likedBy.indexOf(username, 0);
    if (index > -1) {
      const newLikedBy: string[] = []
      post.likedBy.forEach((user, i) => {
        if (index !== i) {
          newLikedBy.push(user);
        }
      });
      post.likedBy = newLikedBy;
      liked = false;
    } else {
      const newLikedBy: string[] = [username];
      post.likedBy.forEach(user => newLikedBy.push(user));
      post.likedBy = newLikedBy;
    }

    this.gqlService.likePost(post, liked);
  }

  editPost(post: Post) {
    post.editMode = false;
    this.gqlService.editPost(post);
  }

  removePost(post: Post) {
    post.editMode = false;
    this.gqlService.removePost(post);
  }

  showComments(post: Post) {
    post.commentMode = true;
    this.gqlService.getPostComments(post);
  }
}
