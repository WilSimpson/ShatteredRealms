import { TestBed } from '@angular/core/testing';

import { UsersService } from './users.service';
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {HttpClient} from "@angular/common/http";
import {of} from "rxjs";

describe('UsersService', () => {
  let service: UsersService;
  let http: HttpClient;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule]
    }, );
    service = TestBed.inject(UsersService);
    http = TestBed.inject(HttpClient);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should be able to get users', () => {
    const resp = {id: 1}
    spyOn(http, 'get').and.returnValue(of(resp));
    service.getUser(resp.id).subscribe((result) => {
      expect(http.get).toHaveBeenCalledTimes(1);
      expect(result).toEqual(resp);
    });
  });


  it('should be able to get users', () => {
    const resp = {id: 1}
    spyOn(http, 'get').and.returnValue(of(resp));
    service.getAllUsers().subscribe((result) => {
      expect(http.get).toHaveBeenCalledTimes(1);
      expect(result).toEqual(resp);
    });
  });
});
