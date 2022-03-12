import {Component, OnInit} from '@angular/core';
import {Post} from "../../model/post";
import {GQLService} from "../../service/gql.service";
import {ActivatedRoute} from "@angular/router";
import {Location} from "@angular/common";

@Component({
  selector: 'app-single-post',
  templateUrl: './single-post.component.html',
  styleUrls: ['./single-post.component.scss'],
  exportAs: 'SinglePostComponent'
})
export class SinglePostComponent implements OnInit {

  post: Post | null;
  valid: boolean = true;

  constructor(private gqlService: GQLService,
                    private route: ActivatedRoute,
                    private location: Location) {
    let postID = this.route.snapshot.paramMap.get('id');
    // get postID when opened via dialog
    if (postID === null) {
      postID = this.location.path().substring(3);
    }

    if (postID !== null) {
      let returnedData = this.gqlService.getPostByID(postID);

      if (returnedData instanceof Post) {
        this.post = returnedData;
      } else {
        returnedData.subscribe({
          next: ({data}: any) => {
            console.log(data)
            this.post = new Post(data.getPost);
          },
          error: (_: any) => {
            this.valid = false;
          }
        });
      }
    } else {
      this.valid = false
    }

  }

  ngOnInit(): void {
  }
}
