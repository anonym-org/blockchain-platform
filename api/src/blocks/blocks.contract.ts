import { ApiProperty } from '@nestjs/swagger';
import { IsNotEmpty, IsObject } from 'class-validator';

export class CreateDto {
  @IsNotEmpty()
  @IsObject()
  @ApiProperty()
  data: any;
}
