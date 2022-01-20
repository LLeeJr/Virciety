import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {Post} from "../../model/post";
import {GQLService} from "../../service/gql.service";
import {AuthLibService} from "auth-lib";

@Component({
  selector: 'post-comment',
  templateUrl: './comment.component.html',
  styleUrls: ['./comment.component.scss']
})
export class CommentComponent implements OnInit {

  @Input() post: Post;
  @Output() newCommentEvent = new EventEmitter<string>()
  comment: string = "";
  nameSourceMap: Map<string, any> = new Map<string, any>();
  source: string = '';

  constructor(private auth: AuthLibService,
              private gqlService: GQLService) {
  }

  ngOnInit(): void {
    this.auth.getUserByName(this.post.username).subscribe(value => {
      if (value && value.data && value.data.getUserByName) {
        let {profilePictureId} = value.data.getUserByName;
        if (profilePictureId && profilePictureId !== '') {
          this.auth.getProfilePicture(profilePictureId).subscribe(picture => {
            if (picture && picture.data && picture.data.getProfilePicture) {
              this.source = picture.data.getProfilePicture;
            }
          })
        }
      }
    });
    this.gqlService.userProfilePictureIds.subscribe(map => {
      for (let entry of map.entries()) {
        if (entry[1] !== "") {
          this.auth.getProfilePicture(entry[1]).subscribe(value => {
            if (value && value.data && value.data.getProfilePicture) {
              this.nameSourceMap.set(entry[0], value.data.getProfilePicture);
            }
          })
        } else {
          this.nameSourceMap.set(entry[0], "");
        }
      }
    })
  }

  addComment() {
    this.newCommentEvent.emit(this.comment);
    this.comment = "";
  }

}
