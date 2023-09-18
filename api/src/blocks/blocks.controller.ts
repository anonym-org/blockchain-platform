import { Body, Controller, Get, Logger, Post } from '@nestjs/common';
import { CreateDto } from './blocks.contract';
import { ApiTags } from '@nestjs/swagger';
import { BlocksService } from './blocks.service';

@ApiTags('Blocks')
@Controller('blocks')
export class BlocksController {
  constructor(private readonly blocksService: BlocksService) {}

  @Get()
  async findAllBlocks(): Promise<any> {
    return await this.blocksService.findAll();
  }

  @Post()
  async createBlocks(@Body() createDto: CreateDto): Promise<any> {}
}
