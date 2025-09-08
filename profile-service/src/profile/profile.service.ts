import { Injectable } from '@nestjs/common';
import { PrismaClient } from '../generated/prisma';
import { CreateProfileDto } from './dto/create-profile.dto';
import { UpdateProfileDto } from './dto/update-profile.dto';
import { Profile } from './profile.entity';

@Injectable()
export class ProfileService {
  private prisma = new PrismaClient();

  async create(createProfileDto: CreateProfileDto): Promise<Profile> {
    return this.prisma.profile.create({
      data: createProfileDto,
    });
  }

  async findAll(): Promise<Profile[]> {
    return this.prisma.profile.findMany();
  }

  async findOne(id: string): Promise<Profile | null> {
    return this.prisma.profile.findUnique({
      where: { id },
    });
  }

  async findByUserId(userId: string): Promise<Profile | null> {
    return this.prisma.profile.findFirst({
      where: { userId },
    });
  }

  async update(id: string, updateProfileDto: UpdateProfileDto): Promise<Profile> {
    return this.prisma.profile.update({
      where: { id },
      data: updateProfileDto,
    });
  }

  async remove(id: string): Promise<Profile> {
    return this.prisma.profile.delete({
      where: { id },
    });
  }
}