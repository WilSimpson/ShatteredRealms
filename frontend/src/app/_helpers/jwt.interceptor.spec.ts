import { TestBed } from '@angular/core/testing';

import { JwtInterceptor } from './jwt.interceptor';
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {AuthenticationService} from "../_services/authentication.service";
import {HttpEvent, HttpHandler, HttpRequest} from "@angular/common/http";
import {hasNameIdentifier} from "@angular/compiler-cli/ngcc/src/utils";
import {Observable} from "rxjs";

describe('JwtInterceptor', () => {
  let jwtInterceptor: JwtInterceptor;
  let authenticationService: AuthenticationService;
  let request: HttpRequest<any>;
  let next: HttpHandler;
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        JwtInterceptor
      ],
      imports: [HttpClientTestingModule],
    }).compileComponents();
    jwtInterceptor = TestBed.inject(JwtInterceptor);
    authenticationService = TestBed.inject(AuthenticationService);
    request = new HttpRequest<any>('GET', 'https://test.com');
    next = new MockHttpHandler();
    spyOn(next, 'handle');
  });

  it('should be created', () => {
    expect(jwtInterceptor).toBeTruthy();
  });

  describe('intercept', () => {
    it('should not set headers if no user is signed in', () => {
      jwtInterceptor.intercept(request, next);
      expect(next.handle).toHaveBeenCalledOnceWith(request);
    });

    it('should set headers if a user is signed in', () => {
      spyOnProperty(authenticationService, 'currentUserValue').and.returnValue({token: 'asdf'});
      jwtInterceptor.intercept(request, next);
      const newRequest = request.clone({
        setHeaders: {
          Authorization: `Bearer asdf`
        }
      });
      expect(next.handle).toHaveBeenCalledWith(newRequest);
    });
  });
});

class MockHttpHandler extends HttpHandler {
  handle(req: HttpRequest<any>): Observable<HttpEvent<any>> {
    return new Observable<HttpEvent<any>>();
  }
}
