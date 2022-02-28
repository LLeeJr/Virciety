import { TestBed } from '@angular/core/testing';

import { GQLService } from './gql.service';

describe('GqlService', () => {
  let service: GQLService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(GQLService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
